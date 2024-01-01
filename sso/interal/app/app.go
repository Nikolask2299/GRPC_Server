package grpcapp

import (
	"log/slog"
	authgrpc "sso/interal/app"
	"time"
	"google.golang.org/grpc"
)

type App struct {
	gRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration,) *App {
	
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
		
	}
}