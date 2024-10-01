package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"goph_keeper/goph_server/internal/config"
	"goph_keeper/goph_server/internal/grpc"
	"goph_keeper/goph_server/internal/logging"
	"goph_keeper/goph_server/internal/storage/postgresql"
)

func Run() error {
	logging.Log().Infof("starting server")
	serverConfig, err := config.ProcessFlags()
	if err != nil {
		return fmt.Errorf("flag processing error: %w", err)
	}
	repo, err := postgresql.NewPostgreSQLRepo(serverConfig)
	if err != nil {
		return fmt.Errorf("postgresql repository error: %w", err)
	}

	serverCtx, serverCtxCancel := context.WithCancel(context.Background())
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		sig := <-sigint
		logging.Log().Infof("received signal: %s", sig)
		serverCtxCancel()
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverConfig.GrpcPort))
	if err != nil {
		return fmt.Errorf("listen port error: %w", err)
	}
	server := grpc.RunGrpcServer(serverCtx, repo)
	if err = server.Serve(listener); err != nil {
		return err
	}

	return nil
}
