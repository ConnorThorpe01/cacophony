package main

import (
	cacophony "cacophony/proto"
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt/v4"
	uuid2 "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"server/auth"
	"server/db"
	"time"
)

type server struct {
	cacophony.UnimplementedChatServiceServer
	db *sql.DB
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
	err = db.CreateAccount(s.db, r.Username, string(hashedPassword), r.Email, uuid)

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
	hash_pass, userID, err := db.Login(s.db, r.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password and/or username")
	}
	if bcrypt.CompareHashAndPassword([]byte(hash_pass), []byte(r.Password)) != nil {
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

	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (s *server) ReceiveMessage(r *cacophony.ReceiveMessageRequest, m grpc.ServerStreamingServer[cacophony.Message]) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveMessage not implemented")
}
func (s *server) CreateGroup(c context.Context, r *cacophony.CreateGroupRequest) (*cacophony.CreateGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateGroup not implemented")
}
func (s *server) AddUserToGroup(c context.Context, r *cacophony.AddUserToGroupRequest) (*cacophony.AddUserToGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddUserToGroup not implemented")
}
func (s *server) SendGroupMessage(c context.Context, r *cacophony.SendGroupMessageRequest) (*cacophony.SendGroupMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendGroupMessage not implemented")
}
func (s *server) ReceiveGroupMessage(r *cacophony.ReceiveGroupMessageRequest, m grpc.ServerStreamingServer[cacophony.GroupMessage]) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveGroupMessage not implemented")
}
