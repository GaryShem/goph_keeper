package grpc

import "goph_keeper/shared/proto"

type Card struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	CVV    string `json:"cvv"`
}

func CardFromProto(protoCard *proto.Card) Card {
	if protoCard == nil {
		return Card{}
	}
	return Card{
		Name:   protoCard.Name,
		Number: protoCard.Number,
		CVV:    protoCard.Cvv,
	}
}
