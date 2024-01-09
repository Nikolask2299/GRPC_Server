package auth

import (
	"context"
	"errors"
	ssov1 "protos/gen/go/sso"
	"sso/interal/services/auth"
	"sso/interal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyId = 0
	emptyEmPw = ""
)


type Auth interface {
	Login(ctx context.Context, email string, password string, appId int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth){
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	
	if err := isValidationLogin(req); err != nil {
		return nil, err
	}

 	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId())) 
	if err != nil {
		if errors.Is(err,  auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "Invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error){
	if err := isValidationRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{UserId: int64(userID),}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error){
	if err := isValidationAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "User not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}







func isValidationLogin(req *ssov1.LoginRequest) error {

	if req.GetEmail() == emptyEmPw {
		return status.Error(codes.InvalidArgument, "Email cannot be empty")
	}

	if req.GetPassword() == emptyEmPw {
        return status.Error(codes.InvalidArgument, "Password cannot be empty")
    }

	if req.GetAppId() == emptyId {
		return status.Error(codes.InvalidArgument, "AppId cannot be empty")
	}	
	
	return nil
}

func isValidationRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == emptyEmPw {
		return status.Error(codes.InvalidArgument, "Email cannot be empty")
	}

	if req.GetPassword() == emptyEmPw {
        return status.Error(codes.InvalidArgument, "Password cannot be empty")
    }

	return nil
}

func isValidationAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyId {
		return status.Error(codes.InvalidArgument, "UserId cannot be empty")
	}

	return nil
}