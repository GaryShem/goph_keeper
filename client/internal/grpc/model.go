package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"goph_keeper/shared/proto"
)

var ErrClientNotReady = errors.New("grpc client not ready, log in first")
var ErrSettingsIncomplete = errors.New("grpc settings incomplete")

var grpcWrapper GrpcWrapper

type KeeperData struct {
	DataName string
	Text     string
	Card     Card
	Binary   []byte
}

type GrpcWrapper struct {
	clientReady bool
	settings    ClientSettings
	grpcClient  proto.GophKeeperClient
}

func (g *GrpcWrapper) InitClient() error {
	if g.clientReady {
		return nil
	}
	if g.settings.Username == "" || g.settings.Password == "" || g.settings.Host == "" || g.settings.Port == 0 {
		return ErrSettingsIncomplete
	}
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", g.settings.Host, g.settings.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	g.grpcClient = proto.NewGophKeeperClient(conn)
	g.clientReady = true
	return nil
}

func Client() *GrpcWrapper {
	return &grpcWrapper
}

func (g *GrpcWrapper) Ping() error {
	if !g.clientReady {
		return ErrClientNotReady
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := g.grpcClient.Ping(ctx, &emptypb.Empty{})
	return err
}
