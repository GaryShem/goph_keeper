package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (g *GophKeeper) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
