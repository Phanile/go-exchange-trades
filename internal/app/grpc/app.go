package grpc

import (
	"fmt"
	"github.com/Phanile/go-exchange-trades/internal/config"
	grpcTrades "github.com/Phanile/go-exchange-trades/internal/grpc/trades"
	"github.com/Phanile/go-exchange-trades/internal/middleware"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewGRPCApp(log *slog.Logger, config *config.GRPCConfig, tradesService grpcTrades.Trades, middleware *middleware.JWTMiddleware) *App {
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.UnaryInterceptor()),
	)

	grpcTrades.Register(gRPCServer, tradesService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       config.Port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const op = "grpcApp.Run"

	app.log.With(
		slog.String("op", op),
		slog.Int("port", app.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))

	if err != nil {
		return fmt.Errorf("%s : %w", op, err)
	}

	app.log.Info("grpc server is running on", slog.String("addr", listener.Addr().String()))

	return app.gRPCServer.Serve(listener)
}

func (app *App) Stop() {
	const op = "grpcApp.Stop"

	app.log.With(
		slog.String("op", op),
	)

	app.log.Info("grpc server is shutting down")

	app.gRPCServer.GracefulStop()
}
