package postgresql

func (s *PostgreSQLSuite) TestRegister() {
	err := s.repo.RegisterUser("test_register", "test_register")
	s.Require().NoError(err)
}

func (s *PostgreSQLSuite) TestLogin() {
	err := s.repo.LoginUser("test", "test")
	s.Require().NoError(err)
	err = s.repo.LoginUser("test", "test2")
	s.Require().Error(err)
}

func (s *PostgreSQLSuite) TestLoginNoUser() {
	err := s.repo.LoginUser("pdfjpsjgpsjf", "test")
	s.Require().ErrorContains(err, "user not found")
}
