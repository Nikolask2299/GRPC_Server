package auth

import (
	"context"
	ssov1 "protos/gen/go/sso"

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
	RegisterNewUser(ctx context.Context, email string, password string) (userID int, err error)
	IsAdmin(ctx context.Context, userID int) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth){
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	
	if req.GetEmail() == emptyEmPw {
		return nil, status.Error(codes.InvalidArgument, "Email cannot be empty")
	}

	if req.GetPassword() == emptyEmPw {
        return nil, status.Error(codes.InvalidArgument, "Password cannot be empty")
    }

	if req.GetAppId() == emptyId {
		return nil, status.Error(codes.InvalidArgument, "AppId cannot be empty")
	}

 	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId())) 
	if err != nil {
		//TODO: log error and return error to client
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error){
	panic("implement me")
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error){
	panic("implement me")
}



func isValidationRequest(req *ssov1.LoginRequest) error {

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