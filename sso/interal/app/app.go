package app

import (
	"log/slog"
	"time"
	"google.golang.org/grpc"
)

type App struct {
	gRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration,) *App {
	
	grpcApp := grpcapp.New() 

	return &App{
		GRPCSrv: grpcApp,

	}
}