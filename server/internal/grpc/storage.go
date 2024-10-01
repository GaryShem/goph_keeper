package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"goph_keeper/shared/proto"
)

func (g *GophKeeper) DownloadStorage(ctx context.Context, _ *emptypb.Empty) (*proto.DownloadStorageResponse, error) {
	user := getUser(ctx)
	data, err := g.repo.DownloadStorage(user)
	if err != nil {
		return nil, err
	}
	result := make([]*proto.DataMessage, len(data))
	for i, d := range data {
		grpcData, mapErr := mapDataInternalToGrpc(&d)
		if mapErr != nil {
			return nil, mapErr
		}
		result[i] = grpcData
	}
	return &proto.DownloadStorageResponse{Data: result}, nil
}
