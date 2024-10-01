package repo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRepoDataSuite(t *testing.T) {
	suite.Run(t, new(RepoDataSuite))
}

type RepoDataSuite struct {
	suite.Suite
	data RepoData
}

func (s *RepoDataSuite) BeforeTest(_, _ string) {
	s.data = RepoData{
		User:   "",
		Name:   "",
		Type:   "",
		Text:   "",
		Card:   Card{},
		Binary: nil,
	}
}

func (s *RepoDataSuite) TestNoUser() {
	err := s.data.Validate(false)
	s.Require().ErrorIs(err, ErrNoUser)
}

func (s *RepoDataSuite) TestNoName() {
	s.data.User = "foo"
	err := s.data.Validate(false)
	s.Require().ErrorIs(err, ErrNoName)
}

func (s *RepoDataSuite) TestNoTypeRead() {
	s.data.User = "foo"
	s.data.Name = "foo"
	err := s.data.Validate(false)
	s.Require().NoError(err)
}

func (s *RepoDataSuite) TestNoTypeWrite() {
	s.data.User = "foo"
	s.data.Name = "foo"
	err := s.data.Validate(true)
	s.Require().ErrorIs(err, ErrInvalidDataType)
}

func (s *RepoDataSuite) TestNoData() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = TextType
	err := s.data.Validate(true)
	s.Require().ErrorIs(err, ErrNoDataToWrite)
}

func (s *RepoDataSuite) TestMultipleData() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = TextType
	s.data.Text = "wololo"
	s.data.Card = Card{Name: "foo"}
	s.data.Binary = []byte{1, 2, 3, 4}
	err := s.data.Validate(true)
	s.Require().ErrorIs(err, ErrHasMultipleData)
}

func (s *RepoDataSuite) TestNoText() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = TextType
	s.data.Card = Card{Name: "foo"}
	err := s.data.Validate(true)
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrNoDataToWrite)
}

func (s *RepoDataSuite) TestInvalidCard() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = CardType
	s.data.Card = Card{
		Name:   "foo bar",
		Number: "1234123412341234",
		CVV:    "12a",
	}
	err := s.data.Validate(true)
	s.Require().ErrorIs(err, ErrInvalidCVV)
}

func (s *RepoDataSuite) TestInvalidBinary() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = BinaryType
	s.data.Card = Card{Name: "foo"}
	err := s.data.Validate(true)
	s.Require().ErrorIs(err, ErrNoDataToWrite)
}

func (s *RepoDataSuite) TestValidText() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = TextType
	s.data.Text = "wololo"
	err := s.data.Validate(true)
	s.Require().NoError(err)
}

func (s *RepoDataSuite) TestValidCard() {
	s.data.User = "foo"
	s.data.Name = "foo"
	s.data.Type = CardType
	s.data.Card = Card{
		Name:   "foo bar",
		Number: "1234123412341234",
		CVV:    "123",
	}
	err := s.data.Validate(true)
	s.Require().NoError(err)
}
