package user

import (
	"context"
	"encoding/json"
	"github.com/bitbeliever/creve-grpc/model"
	"github.com/bitbeliever/creve-grpc/pkg/jwt"
	userservicepb "github.com/bitbeliever/creve-grpc/proto/user/pb"
	jwtutil "github.com/bitbeliever/microutils/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
)

type Service struct {
	userservicepb.UnimplementedUserServiceServer
}

func (s *Service) Hello(ctx context.Context, req *userservicepb.HelloRequest) (*userservicepb.HelloResponse, error) {
	log.Printf("Received: %v", req)
	return &userservicepb.HelloResponse{
		Message: "Hello juz",
	}, nil
}

func (s *Service) NewUser(ctx context.Context, in *userservicepb.NewUserRequest) (*userservicepb.NewUserResponse, error) {
	u, err := model.NewUser(in)
	if err != nil {
		return &userservicepb.NewUserResponse{}, err
	}

	return &userservicepb.NewUserResponse{
		Username: u.Username,
		Token:    u.Password,
	}, nil
}

func (s *Service) GetUserByName(ctx context.Context, in *userservicepb.GetUserByNameRequest) (*userservicepb.UserResponse, error) {
	log.Printf("in %v", in)
	u, err := model.GetUserByName(in.Username)
	if err != nil {
		return &userservicepb.UserResponse{}, status.New(codes.NotFound, "User not found").Err()
	}

	return &userservicepb.UserResponse{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Token:    "todo",
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, in *userservicepb.GetUserByIDRequest) (*userservicepb.UserResponse, error) {
	u, err := model.GetUserByID(in.Id)
	if err != nil {
		return &userservicepb.UserResponse{}, err
	}

	return &userservicepb.UserResponse{
		Username: u.Username,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Token:    "todo",
		Id:       u.ID,
	}, nil
}

func (s *Service) Login(ctx context.Context, in *userservicepb.LoginRequest) (*userservicepb.UserResponse, error) {
	u, err := model.GetUserByName(in.Username)
	if err != nil {
		return &userservicepb.UserResponse{}, status.New(codes.NotFound, err.Error()).Err()
	}

	if u.Password != in.Password {
		return &userservicepb.UserResponse{}, status.New(20001, "incorrect pwd").Err()
	}

	token, err := jwt.NewToken(u)
	if err != nil {
		return &userservicepb.UserResponse{}, status.New(codes.Unknown, "sign token error").Err()
	}
	return &userservicepb.UserResponse{
		Username: u.Username,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Token:    token,
		Id:       u.ID,
		UserId:   u.UserID,
	}, nil
}

// Test test query string
func (s *Service) Test(ctx context.Context, in *userservicepb.TestRequest) (*userservicepb.TestResponse, error) {
	log.Printf("in %v", in)
	if in.Query != "query" {
		return &userservicepb.TestResponse{}, status.New(codes.InvalidArgument, "invalid query").Err()
	}
	m, err := jwt.ParseToken(in.Token)
	if err != nil {
		log.Println(err)
	}

	log.Println(jwtutil.IsExpire(in.Token))

	_ = json.NewEncoder(os.Stdout).Encode(m)
	return &userservicepb.TestResponse{
		Name: in.Name + "test",
		Age:  in.Age + "test  af",
	}, nil
}
