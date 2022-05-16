package main

import (
	"context"

	pb "github.com/jimdhughes/go-user-ms/proto"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) RegisterUser(ctx context.Context, in *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	u := User{
		Email:    in.Email,
		Password: in.Password,
	}
	success, err := DB.CreateUser(u)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterUserResponse{
		Success: success,
	}, nil

}

func (s *server) LoginUser(ctx context.Context, in *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	token, err := DB.Login(in.Email, in.Password)
	return &pb.LoginUserResponse{
		Token: token,
	}, err
}

func (s *server) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	userSafe, err := TS.ValidateToken(in.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, err
	}
	return &pb.ValidateTokenResponse{
		Valid: true,
		ID:    userSafe.ID,
		Email: userSafe.Email,
	}, nil

}
