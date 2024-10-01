package grpc

import (
	"context"
	"time"

	"goph_keeper/shared/proto"
)

func (g *GrpcWrapper) RegisterUser() error {
	if !g.clientReady {
		return ErrClientNotReady
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := g.grpcClient.RegisterUser(ctx, &proto.RegisterRequest{
		Name:     g.settings.Username,
		Password: g.settings.Password,
	})
	return err
}
