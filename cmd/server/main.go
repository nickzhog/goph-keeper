package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/nickzhog/goph-keeper/internal/server/config"
	"github.com/nickzhog/goph-keeper/internal/server/server"
	"github.com/nickzhog/goph-keeper/internal/server/server/grpc"
	"github.com/nickzhog/goph-keeper/internal/server/service"
	"github.com/nickzhog/goph-keeper/pkg/logging"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()
	logger.Tracef("config: %+v", cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		oscall := <-c
		logger.Tracef("system call:%+v", oscall)
		cancel()
	}()

	storage := service.NewPostgresStorage(ctx, logger, cfg.PostgresStorage.DatabaseDSN)

	srv := server.NewServer(logger, cfg, storage)

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		grpc.Serve(ctx, *srv, cfg)
		wg.Done()
	}()

	wg.Wait()
}
