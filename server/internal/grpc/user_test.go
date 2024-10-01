package grpc

import "goph_keeper/shared/proto"

func (s *GrpcSuite) TestRegister() {
	ctx, cancel := s.getContext()
	defer cancel()
	_, err := s.agent.RegisterUser(ctx, &proto.RegisterRequest{
		Name:     "grpc_test",
		Password: "test",
	})
	s.Require().NoError(err)
}

func (s *GrpcSuite) TestDuplicateRegister() {
	ctx, cancel := s.getContext()
	defer cancel()
	_, err := s.agent.RegisterUser(ctx, &proto.RegisterRequest{
		Name:     "grpc_duplicate",
		Password: "test",
	})
	s.Require().NoError(err)
	_, err = s.agent.RegisterUser(ctx, &proto.RegisterRequest{
		Name:     "grpc_duplicate",
		Password: "test",
	})
	s.Require().Error(err)
}

func (s *GrpcSuite) TestNoAuth() {
	ctx, cancel := s.getContext()
	defer cancel()
	_, err := s.agent.SetData(ctx, nil)
	s.Require().Error(err)
	s.Require().ErrorContains(err, "access denied")
}
