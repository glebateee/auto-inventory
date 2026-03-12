package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/glebateee/auto-inventory/internal/grpc/server"
	"google.golang.org/grpc"
)

type App struct {
	logger *slog.Logger
	addr   string
	srv    *grpc.Server
}

func New(
	logger *slog.Logger,
	host string,
	port int,
	provider server.Provider,
) *App {
	srv := grpc.NewServer()
	server.Register(srv, provider)
	return &App{
		logger: logger,
		srv:    srv,
		addr:   net.JoinHostPort(host, strconv.Itoa(port)),
	}
}

func (a *App) Start() error {
	const op = "grpcapp.Start"
	logger := a.logger.With(
		slog.String("op", op),
	)
	lis, err := net.Listen("tcp", a.addr)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	logger.Info("server started", slog.String("addr", a.addr))
	if err := a.srv.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	logger := a.logger.With(
		slog.String("op", op),
	)
	a.srv.GracefulStop()
	logger.Info("server stopped")
}
