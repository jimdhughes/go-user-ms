package main

import (
	"context"
	"errors"
	"fmt"
	"log"

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
	if err != nil {
		log.Println(err)
		return nil, errors.New("invalid username/password combination")
	}
	return &pb.LoginUserResponse{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, err
}

func (s *server) ValidateAccessToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if in.AccessToken == "" {
		return nil, errors.New("no access token provided")
	}
	userSafe, err := TS.ValidateAccessToken(in.AccessToken)
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

func (s *server) RefreshAccessToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.RefreshAccessTokenResponse, error) {
	if in.RefreshToken == "" {
		return nil, errors.New("no refresh token provided")
	}
	userId, err := TS.ValidateRefreshToken(in.RefreshToken)
	if err != nil {
		return nil, err
	}
	if userId == "" {
		return nil, fmt.Errorf("invalid refresh token")
	}
	user, err := DB.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	tokenPair, err := TS.GenerateTokenPairForUser(user)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshAccessTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil
}
