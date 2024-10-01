package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	"goph_keeper/shared/proto"
)

func (g *GrpcWrapper) SetData(dataType proto.DATA_TYPE, data KeeperData) error {
	if !g.clientReady {
		return ErrClientNotReady
	}
	md := metadata.Pairs("username", g.settings.Username, "password", g.settings.Password)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	request := &proto.DataMessage{
		Name:     data.DataName,
		DataType: dataType,
		Text:     data.Text,
		Card: &proto.Card{
			Name:   data.Card.Name,
			Number: data.Card.Number,
			Cvv:    data.Card.CVV,
		},
		Binary: data.Binary,
	}
	_, err := g.grpcClient.SetData(ctx, request)
	return err
}

func (g *GrpcWrapper) GetData(dataname string) (*KeeperData, error) {
	var result *KeeperData
	if !g.clientReady {
		return nil, ErrClientNotReady
	}
	md := metadata.Pairs("username", g.settings.Username, "password", g.settings.Password)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	grpcData, err := g.grpcClient.GetData(ctx, &proto.GetDataRequest{
		DataName: dataname,
	})
	if err != nil {
		return nil, err
	}
	result = &KeeperData{
		DataName: grpcData.Name,
		Text:     grpcData.Text,
		Card:     CardFromProto(grpcData.Card),
		Binary:   grpcData.Binary,
	}
	return result, err
}
