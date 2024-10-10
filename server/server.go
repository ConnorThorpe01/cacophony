package main

import (
	cacophony "cacophony/proto"
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	uuid2 "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"server/auth"
	"server/db"
	"strconv"
	"time"
)

type server struct {
	cacophony.UnimplementedChatServiceServer
	sqlDB *sql.DB
	rDB   *redis.Client
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
func (s *server) SendMessage(c context.Context, r *cacophony.SendMessageRequest) (*cacophony.SendMessageResponse, error) {
	userID, err := auth.ValidateJWT(r.Token, jwtSecret)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "error validating JWT token")
	}
	r.Message.FromUserId = userID
	_, err = db.StoreMessage(s.sqlDB, r.Message, r.ConversationID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error storing message")
	}

	ctx := context.Background()
	err = s.rDB.Publish(ctx, strconv.FormatUint(r.ConversationID, 10), r.Message).Err()
	if err != nil {
		return nil, err
	}

	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}

func (s *server) ReceiveMessage(r *cacophony.ReceiveMessageRequest, stream cacophony.ChatService_ReceiveMessageServer) error {
	userID, err := auth.ValidateJWT(r.Token, jwtSecret) // Parse token for user ID
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "error validating JWT token")
	}
	// Example: Wait for new messages from Redis or another system for this user
	pubsub := s.rDB.Subscribe(stream.Context(), userID)
	defer pubsub.Close()

	// Continuously receive messages and push them to the stream
	for {
		select {
		case <-stream.Context().Done():
			return nil // Client closed the stream
		case msg := <-pubsub.Channel():
			// Send the message to the client stream
			err := stream.Send(&cacophony.Message{
				FromUserId: "otherUserId", // Use actual sender ID
				Content:    msg.Payload,
				Timestamp:  time.Now().Unix(),
			})
			if err != nil {
				return err
			}
		}
	}
}

func (s *server) CreateGroup(c context.Context, r *cacophony.CreateGroupRequest) (*cacophony.CreateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (s *server) AddUserToGroup(c context.Context, r *cacophony.AddUserToGroupRequest) (*cacophony.AddUserToGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserToGroup not implemented")
}
