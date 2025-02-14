package auth

import (
	"context"

	ssov1 "github.com/MrFlameZahar/grpctest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}
func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{UserId: userID}, nil
}
func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}
	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, err
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func validateLogin(req *ssov1.LoginRequest) error {
	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}
	if req.GetAppId() == 0 {
		return status.Error(codes.InvalidArgument, "invalid appID")
	}

	return nil
}

func validateRegister(req *ssov1.RegisterRequest) error {
	if req.Email == "" {
		return status.Error(codes.InvalidArgument, "email is empty")
	}
	if req.Password == "" {
		return status.Error(codes.InvalidArgument, "password is empty")
	}

	return nil
}

func validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == 0 {
		return status.Error(codes.InvalidArgument, "invalid userID")
	}

	return nil
}
