package grpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"

	interceptors2 "goph_keeper/goph_server/internal/interceptors"
	"goph_keeper/goph_server/internal/logging"
	"goph_keeper/goph_server/internal/storage/repo"
	"goph_keeper/shared/proto"
)

type GophKeeper struct {
	proto.UnimplementedGophKeeperServer
	repo repo.Repo
}

func NewGophKeeper(repo repo.Repo) *GophKeeper {
	return &GophKeeper{repo: repo}
}

func RunGrpcServer(ctx context.Context, repo repo.Repo) *grpc.Server {
	keeper := NewGophKeeper(repo)
	interceptors := []grpc.UnaryServerInterceptor{
		interceptors2.MethodLogInterceptor{}.Intercept,
		interceptors2.NewAuthInterceptor(repo).Intercept,
	}
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))
	proto.RegisterGophKeeperServer(server, keeper)

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownStopCtx := context.WithTimeout(context.Background(), 10*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				server.Stop()
				logging.Log().Warn("server could not stop gracefully, forcing shutdown")
			}
		}()
		logging.Log().Info("trying to gracefully shutdown server")
		server.GracefulStop()
		logging.Log().Info("server shut down")
		shutdownStopCtx()
	}()
	return server
}

var _ proto.GophKeeperServer = (*GophKeeper)(nil)
