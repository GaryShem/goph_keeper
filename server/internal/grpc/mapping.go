package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"

	"goph_keeper/goph_server/internal/storage/repo"
	"goph_keeper/shared/proto"
)

func getUser(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["username"]) == 0 {
		return ""
	}
	return md["username"][0]
}

func mapCardGrpcToInternal(card *proto.Card) repo.Card {
	if card == nil {
		return repo.Card{}
	}
	result := repo.Card{
		Name:   card.Name,
		Number: card.Number,
		CVV:    card.Cvv,
	}
	return result
}

func mapCardInternalToGrpc(card *repo.Card) *proto.Card {
	if card == nil || !card.HasData() {
		return nil
	}
	return &proto.Card{
		Name:   card.Name,
		Number: card.Number,
		Cvv:    card.CVV,
	}
}

func mapDataTypeGrpcToInternal(dataType proto.DATA_TYPE) repo.DataType {
	switch dataType {
	case proto.DATA_TYPE_text:
		return repo.TextType
	case proto.DATA_TYPE_binary:
		return repo.BinaryType
	case proto.DATA_TYPE_card:
		return repo.CardType
	default:
		return ""
	}
}

func mapDataTypeInternalToGrpc(dataType repo.DataType) (proto.DATA_TYPE, error) {
	switch dataType {
	case repo.TextType:
		return proto.DATA_TYPE_text, nil
	case repo.BinaryType:
		return proto.DATA_TYPE_binary, nil
	case repo.CardType:
		return proto.DATA_TYPE_card, nil
	default:
		return 0, repo.ErrInvalidDataType
	}
}

func mapDataGrpcToInternal(ctx context.Context, data *proto.DataMessage) repo.RepoData {
	return repo.RepoData{
		User:   getUser(ctx),
		Name:   data.Name,
		Type:   mapDataTypeGrpcToInternal(data.DataType),
		Text:   data.Text,
		Card:   mapCardGrpcToInternal(data.Card),
		Binary: data.Binary,
	}
}

func mapDataInternalToGrpc(data *repo.RepoData) (*proto.DataMessage, error) {
	dataType, err := mapDataTypeInternalToGrpc(data.Type)
	if err != nil {
		return nil, err
	}
	return &proto.DataMessage{
		Name:     data.Name,
		DataType: dataType,
		Text:     data.Text,
		Card:     mapCardInternalToGrpc(&data.Card),
		Binary:   data.Binary,
	}, nil
}

func emptyDataGrpcToInternal(ctx context.Context, request *proto.GetDataRequest) repo.RepoData {
	user := getUser(ctx)
	return repo.RepoData{
		User:   user,
		Name:   request.DataName,
		Type:   "",
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	}
}
