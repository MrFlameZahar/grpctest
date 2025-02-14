package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MrFlameZahar/grpctest/internal/app"
	"github.com/MrFlameZahar/grpctest/internal/config"
)

const (
	envLocal = "local"
	envProd  = "Prod"
	envDev   = "Dev"
)

func main() {
	// TODO: object config
	cfg := config.MustLoad()
	// TODO: initialize logger
	log := setupLogger(cfg.Env)

	log.Info("starting app", slog.String("env", cfg.Env))
	// TODO: initialize app
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.Run()

	// TODO: start gRPC-server
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
	return log
}
