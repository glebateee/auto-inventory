package app

import (
	"log/slog"

	"github.com/glebateee/auto-inventory/internal/app/grpcapp"
)

type App struct {
	logger  *slog.Logger
	GRPCApp *grpcapp.App
}

func New(
	logger *slog.Logger,
	host string,
	port int,
) *App {
	gRPCapp := grpcapp.New(
		logger,
		host,
		port,
	)
	return &App{
		GRPCApp: gRPCapp,
	}
}

func (a *App) MustStart() {
	if err := a.GRPCApp.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	a.GRPCApp.Stop()
}
