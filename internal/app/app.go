package app

import (
	"log/slog"

	"github.com/glebateee/auto-inventory/internal/app/grpcapp"
	"github.com/glebateee/auto-inventory/internal/services/provider"
	"github.com/glebateee/auto-inventory/internal/storage/postgres"
)

type App struct {
	logger  *slog.Logger
	GRPCApp *grpcapp.App
}

func New(
	logger *slog.Logger,
	gRPChost string,
	gRPCport int,
	dbname string,
	dbuser string,
	dbpassword string,
	dbhost string,
	dbport int,
	dbsslmode string,
) *App {
	prodProviderStorage, err := postgres.New(
		dbname,
		dbuser,
		dbpassword,
		dbhost,
		dbport,
		dbsslmode,
	)
	if err != nil {
		panic(err)
	}
	prodProviderService := provider.New(logger, prodProviderStorage)
	gRPCapp := grpcapp.New(
		logger,
		gRPChost,
		gRPCport,
		prodProviderService,
	)
	return &App{
		logger:  logger,
		GRPCApp: gRPCapp,
	}
}

func (a *App) MustStart() {
	const op = "app.MustStart"
	logger := a.logger.With(
		slog.String("op", op),
	)
	logger.Info("starting main application dependencies")
	if err := a.GRPCApp.Start(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "app.Stop"
	logger := a.logger.With(
		slog.String("op", op),
	)
	logger.Info("stopping main application dependencies")
	a.GRPCApp.Stop()
}
