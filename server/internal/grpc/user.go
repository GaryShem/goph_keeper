package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"goph_keeper/shared/proto"
)

func (g *GophKeeper) RegisterUser(_ context.Context, request *proto.RegisterRequest) (*emptypb.Empty, error) {
	if err := g.repo.RegisterUser(request.Name, request.Password); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
