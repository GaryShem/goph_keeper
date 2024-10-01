package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"goph_keeper/shared/proto"
)

func (g *GophKeeper) GetData(ctx context.Context, request *proto.GetDataRequest) (*proto.DataMessage, error) {
	data := emptyDataGrpcToInternal(ctx, request)
	if err := data.Validate(false); err != nil {
		return nil, err
	}
	result, err := g.repo.GetData(data)
	if err != nil {
		return nil, err
	}
	err = result.Validate(true)
	if err != nil {
		return nil, err
	}
	response, err := mapDataInternalToGrpc(result)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (g *GophKeeper) SetData(ctx context.Context, request *proto.DataMessage) (*emptypb.Empty, error) {
	data := mapDataGrpcToInternal(ctx, request)
	if err := data.Validate(true); err != nil {
		return nil, err
	}
	err := g.repo.SetData(data)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
