package main

import (
	cacophony "cacophony/proto"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	uuid2 "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
	"server/auth"
	"server/db"
	"strings"
	"sync"
	"time"
)

type server struct {
	cacophony.UnimplementedChatServiceServer
	sqlDB   *sql.DB
	rDB     *redis.Client
	clients sync.Map
}

var jwtSecret = []byte(os.Getenv("JWT_KEY"))

func (s *server) CreateAccount(c context.Context, r *cacophony.CreateAccountRequest) (*cacophony.CreateAccountResponse, error) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "error making new user id")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "error ecripting password")
	}
	err = db.CreateAccount(s.sqlDB, r.Username, string(hashedPassword), r.Email, uuid)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating account in database")
	}

	claims := jwt.MapClaims{
		"user_id": uuid.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating JWT token")
	}
	response := &cacophony.CreateAccountResponse{
		Token:   signedToken,
		Message: "Success",
		Success: true,
	}
	return response, nil
}

func (s *server) Login(c context.Context, r *cacophony.LoginRequest) (*cacophony.LoginResponse, error) {
	userID, hashPass, err := db.Login(s.sqlDB, r.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password and/or username")
	}
	if bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(r.Password)) != nil {
		return nil, status.Errorf(codes.PermissionDenied, "incorrect password and/or username")
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	// Generate the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "error generating JWT token")
	}

	response := &cacophony.LoginResponse{
		Token:   signedToken,
		Success: true,
	}
	return response, nil
}

func (s *server) ChatStream(stream cacophony.ChatService_ChatStreamServer) error {
	// Extract JWT from the stream's context and get the user ID
	token, err := extractJWT(stream.Context())
	if err != nil {
		log.Printf("Error extracting JWT: %v", err)
		return err
	}

	clientUserId, err := auth.ValidateJWT(token, jwtSecret)
	if err != nil {
		log.Printf("Invalid JWT: %v", err)
		return err
	}

	log.Printf("Client %v connected", clientUserId)
	s.clients.Store(clientUserId, stream)

	go s.subscribeToRedisChannel(fmt.Sprintf("user:%s", clientUserId), stream)

	groupIds := s.getUserGroupIds(clientUserId)
	for _, groupId := range groupIds {
		go s.subscribeToRedisChannel(fmt.Sprintf("group:%s", groupId), stream)
	}

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Error receiving message from client %v: %v", clientUserId, err)
				s.clients.Delete(clientUserId)
				return
			}

			s.publishMessageToRedis(msg)
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			log.Printf("Client %v disconnected", clientUserId)
			s.clients.Delete(clientUserId)
			return nil
		}
	}
}

func (s *server) getUserGroupIds(id string) []string {
	groups, err := db.GetGroup(s.sqlDB, id)
	if err != nil {
		log.Printf("Error:%s encountered while geting user groups", err)
		return nil
	}
	return groups

}

func (s *server) CreateGroup(c context.Context, r *cacophony.CreateGroupRequest) (*cacophony.CreateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (s *server) AddUserToGroup(c context.Context, r *cacophony.AddUserToGroupRequest) (*cacophony.AddUserToGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserToGroup not implemented")
}

func (s *server) publishMessageToRedis(msg *cacophony.Message) {
	redisChannel := ""
	if msg.ReceiverType == cacophony.ReceiverType_USER {
		redisChannel = fmt.Sprintf("user:%s", msg.Receiver)
	} else if msg.ReceiverType == cacophony.ReceiverType_GROUP {
		redisChannel = fmt.Sprintf("group:%s", msg.Receiver)
	}

	// Publish message to the appropriate Redis channel
	serializedMessage, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Error serializing message from channel: %v", err)
		return
	}
	err = s.rDB.Publish(context.Background(), redisChannel, string(serializedMessage)).Err()
	if err != nil {
		log.Printf("Error publishing message to Redis: %v", err)
	}
}

func (s *server) subscribeToRedisChannel(channel string, stream cacophony.ChatService_ChatStreamServer) {
	pubsub := s.rDB.Subscribe(context.Background(), channel)
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Printf("Error receiving message from Redis channel %s: %v", channel, err)
			return
		}

		var message cacophony.Message
		err = proto.Unmarshal([]byte(msg.Payload), &message)
		if err != nil {
			log.Printf("Error parsing Redis message: %v", err)
			continue
		}

		err = stream.Send(&message)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			return
		}
	}
}

func (s *server) Friend(c context.Context, r *cacophony.FriendRequest) (*cacophony.FriendResponse, error) {
	token, err := extractJWT(c)
	if err != nil {
		return nil, errors.New("failed to extract user token from context")
	}
	claims, ok := extractClaims(token)
	if !ok {
		return nil, errors.New("failed to extract claims from token")
	}
	clientId, ok := claims["user_id"]
	userId, err := db.Friend(s.sqlDB, r.Username, clientId)

	return nil, nil
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecretString := os.Getenv("JWT_KEY")
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

// Extract JWT from the gRPC metadata
func extractJWT(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("no metadata found in context")
	}

	// Extract the "authorization" field from metadata
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return "", fmt.Errorf("no authorization header provided")
	}

	// JWT token should be in the form "Bearer <token>", so split it
	token := strings.TrimPrefix(authHeader[0], "Bearer ")
	if token == authHeader[0] {
		return "", fmt.Errorf("invalid authorization token format")
	}

	return token, nil
}
