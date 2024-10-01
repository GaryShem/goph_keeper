package grpc

import (
	"sort"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"goph_keeper/shared/proto"
)

func (s *GrpcSuite) TestDownloadStorage() {
	user := "storage_user"
	ctx, cancel := s.getContext()
	defer cancel()
	_, _ = s.agent.RegisterUser(ctx, &proto.RegisterRequest{Name: user, Password: "test"})
	canon := []*proto.DataMessage{
		{
			Name:     "storage_text",
			DataType: proto.DATA_TYPE_text,
			Text:     "wololo",
			Card:     nil,
			Binary:   nil,
		},
		{
			Name:     "storage_card",
			DataType: proto.DATA_TYPE_card,
			Text:     "",
			Card: &proto.Card{
				Name:   "foo bar",
				Number: "1234 1234 1234 1234",
				Cvv:    "123",
			},
			Binary: nil,
		},
		{
			Name:     "storage_binary",
			DataType: proto.DATA_TYPE_binary,
			Text:     "",
			Card:     nil,
			Binary:   []byte{1, 2, 3},
		},
	}
	md := metadata.Pairs("username", user, "password", "test")
	ctx = metadata.NewOutgoingContext(ctx, md)
	for _, m := range canon {
		_, err := s.agent.SetData(ctx, m)
		s.Require().NoError(err)
	}

	resp, err := s.agent.DownloadStorage(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	sort.Slice(canon, func(i, j int) bool {
		return canon[i].Name < canon[j].Name
	})
	sort.Slice(resp.Data, func(i, j int) bool {
		return resp.Data[i].Name < resp.Data[j].Name
	})
	s.Require().Equal(len(canon), len(resp.Data))
	for i := range canon {
		s.Require().EqualExportedValues(canon[i], resp.Data[i])
	}
}
