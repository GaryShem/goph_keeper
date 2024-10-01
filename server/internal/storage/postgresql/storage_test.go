package postgresql

import (
	"sort"

	"goph_keeper/goph_server/internal/storage/repo"
)

func (s *PostgreSQLSuite) TestDownloadStorage() {
	user := "storage_user"
	err := s.repo.RegisterUser(user, "test")
	canon := []repo.RepoData{
		{
			User:   user,
			Name:   "storage_text",
			Type:   repo.TextType,
			Text:   "wololo",
			Card:   repo.Card{},
			Binary: nil,
		},
		{
			User: user,
			Name: "storage_card",
			Type: repo.CardType,
			Text: "",
			Card: repo.Card{
				Name:   "foo bar",
				Number: "1234 1234 1234 1234",
				CVV:    "123",
			},
			Binary: nil,
		},
		{
			User:   user,
			Name:   "storage_binary",
			Type:   repo.BinaryType,
			Text:   "",
			Card:   repo.Card{},
			Binary: []byte{1, 2, 3},
		},
	}

	for _, e := range canon {
		err := s.repo.SetData(e)
		s.Require().NoError(err)
	}

	actual, err := s.repo.DownloadStorage(user)
	s.Require().NoError(err)

	sort.Slice(canon, func(i, j int) bool {
		return canon[i].Name < canon[j].Name
	})
	sort.Slice(actual, func(i, j int) bool {
		return actual[i].Name < actual[j].Name
	})

	s.Require().Equal(actual, canon)
}
