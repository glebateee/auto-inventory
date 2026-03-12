package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/glebateee/auto-inventory/internal/app"
	"github.com/glebateee/auto-inventory/internal/config"
)

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	logger := setupLogger(cfg.Env)
	logger.Debug("loaded configuration", slog.Any("cfg", cfg))
	mainApp := app.New(
		logger,
		cfg.GRPCConfig.Host,
		cfg.GRPCConfig.Port,
		cfg.DBConfig.Name,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.SslMode,
	)
	go mainApp.MustStart()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	logger.Info("received system signal", slog.String("signal", sig.String()))
	mainApp.Stop()
	logger.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
