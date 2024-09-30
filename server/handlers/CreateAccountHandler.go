package handlers

import (
	cacophony "cacophony/proto"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CreateAccountHandler(c context.Context, r *cacophony.CreateAccountRequest) (*cacophony.CreateAccountResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAccount not implemented")
}
