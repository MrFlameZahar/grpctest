package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/MrFlameZahar/grpctest/internal/app/grpc"
	"github.com/MrFlameZahar/grpctest/internal/services/auth"
	"github.com/MrFlameZahar/grpctest/internal/storage/sqlite"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	//TODO: initialize storage
	userStorage, _ := sqlite.New(storagePath)
	//TODO: init auth service
	authService := auth.New(log, userStorage, tokenTTL)

	grpcApp := grpcapp.New(log, grpcPort, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
