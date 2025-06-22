package app

import (
	"github.com/Phanile/go-exchange-trades/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCApp grpc.App
}

func NewApp(port uint, log *slog.Logger) *App {
	gRPCApp := grpc.NewGRPCApp(port, log)

	return &App{
		GRPCApp: *gRPCApp,
	}
}
