package postgresql

import "goph_keeper/goph_server/internal/storage/repo"

func (s *PostgreSQLSuite) TestText() {
	name := "test_text"
	dataType := repo.TextType
	canon := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "wololo",
		Card:   repo.Card{},
		Binary: nil,
	}

	err := s.repo.SetData(canon)
	s.Require().NoError(err)

	actual, err := s.repo.GetData(repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	})
	s.Require().NoError(err)
	s.Require().Equal(*actual, canon)
}

func (s *PostgreSQLSuite) TestCard() {
	name := "test_card"
	dataType := repo.CardType
	canon := repo.RepoData{
		User: "test",
		Name: name,
		Type: dataType,
		Text: "",
		Card: repo.Card{
			Name:   "foo bar",
			Number: "1234 1234 1234 1234",
			CVV:    "123",
		},
		Binary: nil,
	}

	err := s.repo.SetData(canon)
	s.Require().NoError(err)

	actual, err := s.repo.GetData(repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	})
	s.Require().NoError(err)
	s.Require().Equal(*actual, canon)
}

func (s *PostgreSQLSuite) TestBinary() {
	name := "test_binary"
	dataType := repo.BinaryType
	canon := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: []byte{1, 2, 3},
	}

	err := s.repo.SetData(canon)
	s.Require().NoError(err)

	actual, err := s.repo.GetData(repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	})
	s.Require().NoError(err)
	s.Require().Equal(*actual, canon)
}

func (s *PostgreSQLSuite) TestInvalid() {
	name := "test_invalid"
	dataType := repo.TextType
	canon := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	}

	_, err := s.repo.GetData(canon)
	s.Require().Error(err)
}

func (s *PostgreSQLSuite) TestOverwrite() {
	name := "test_overwrite"
	dataType := repo.TextType
	canon := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "wololo",
		Card:   repo.Card{},
		Binary: nil,
	}
	overwrite := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "such wow",
		Card:   repo.Card{},
		Binary: nil,
	}
	err := s.repo.SetData(canon)
	s.Require().NoError(err)
	err = s.repo.SetData(overwrite)
	s.Require().NoError(err)

	actual, err := s.repo.GetData(repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	})
	s.Require().NoError(err)
	s.Require().Equal(*actual, overwrite)
}

func (s *PostgreSQLSuite) TestOverwriteDifferentType() {
	name := "test_overwrite"
	dataType := repo.TextType
	overwriteType := repo.BinaryType
	canon := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "wololo",
		Card:   repo.Card{},
		Binary: nil,
	}
	overwrite := repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   overwriteType,
		Text:   "",
		Card:   repo.Card{},
		Binary: []byte{1, 2, 3},
	}
	err := s.repo.SetData(canon)
	s.Require().NoError(err)
	err = s.repo.SetData(overwrite)
	s.Require().Error(err)

	actual, err := s.repo.GetData(repo.RepoData{
		User:   "test",
		Name:   name,
		Type:   dataType,
		Text:   "",
		Card:   repo.Card{},
		Binary: nil,
	})
	s.Require().NoError(err)
	s.Require().Equal(*actual, canon)
}
