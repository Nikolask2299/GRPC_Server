package app

import (
	"log/slog"
	grpcapp "sso/interal/app/grpc"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration,) *App {
	
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{GRPCSrv: grpcApp,}
}