package main

import (
	cacophony "cacophony/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	cacophony.UnimplementedChatServiceServer
}

func (s *server) CreateAccount(c context.Context, r *cacophony.CreateAccountRequest) (*cacophony.CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
func (s *server) Login(c context.Context, r *cacophony.LoginRequest) (*cacophony.LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (s *server) SendMessage(c context.Context, r *cacophony.SendMessageRequest) (*cacophony.SendMessageResponse, error) {
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
