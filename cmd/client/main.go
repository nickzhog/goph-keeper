package main

import (
	"context"

	"github.com/nickzhog/goph-keeper/internal/client/api/grpc"
	"github.com/nickzhog/goph-keeper/internal/client/cli"
	"github.com/nickzhog/goph-keeper/internal/client/config"
	"github.com/nickzhog/goph-keeper/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetOnlyFileLogger()
	logger.Tracef("config: %+v", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	keeperClient := grpc.NewClient(cfg.Settings.ConnAddress, cfg.Settings.CertSSL, logger)

	cli.Start(ctx, keeperClient, logger)
}
