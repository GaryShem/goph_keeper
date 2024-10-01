package grpc

import (
	"google.golang.org/grpc/metadata"

	"goph_keeper/shared/proto"
)

func (s *GrpcSuite) TestText() {
	name := "grpc_text"
	dataType := proto.DATA_TYPE_text
	ctx, cancel := s.getContext()
	defer cancel()
	canon := proto.DataMessage{
		Name:     name,
		DataType: dataType,
		Text:     "wololo",
		Card:     nil,
		Binary:   nil,
	}
	md := metadata.Pairs("username", "test", "password", "test")
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := s.agent.SetData(ctx, &canon)
	s.Require().NoError(err)

	data, err := s.agent.GetData(ctx, &proto.GetDataRequest{
		DataName: name,
	})
	s.Require().NoError(err)
	s.Require().EqualExportedValues(&canon, data)
}

func (s *GrpcSuite) TestCard() {
	name := "grpc_card"
	dataType := proto.DATA_TYPE_card
	ctx, cancel := s.getContext()
	defer cancel()
	canon := proto.DataMessage{
		Name:     name,
		DataType: dataType,
		Text:     "",
		Card: &proto.Card{
			Name:   "foo bar",
			Number: "1234 1234 1234 1234",
			Cvv:    "123",
		},
		Binary: nil,
	}
	md := metadata.Pairs("username", "test", "password", "test")
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := s.agent.SetData(ctx, &canon)
	s.Require().NoError(err)

	data, err := s.agent.GetData(ctx, &proto.GetDataRequest{
		DataName: name,
	})
	s.Require().NoError(err)
	s.Require().EqualExportedValues(&canon, data)
}

func (s *GrpcSuite) TestBinary() {
	name := "grpc_binary"
	dataType := proto.DATA_TYPE_binary
	ctx, cancel := s.getContext()
	defer cancel()
	canon := proto.DataMessage{
		Name:     name,
		DataType: dataType,
		Text:     "",
		Card:     nil,
		Binary:   []byte{1, 2, 3},
	}
	md := metadata.Pairs("username", "test", "password", "test")
	ctx = metadata.NewOutgoingContext(ctx, md)
	_, err := s.agent.SetData(ctx, &canon)
	s.Require().NoError(err)

	data, err := s.agent.GetData(ctx, &proto.GetDataRequest{
		DataName: name,
	})
	s.Require().NoError(err)
	s.Require().EqualExportedValues(&canon, data)
}
