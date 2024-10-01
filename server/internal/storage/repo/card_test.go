package repo

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCardSuite(t *testing.T) {
	suite.Run(t, new(CardSuite))
}

type CardSuite struct {
	suite.Suite
}

func (s *CardSuite) TestValid() {
	card := Card{
		Name:   "foo bar",
		Number: "12345678 1234  5678    ",
		CVV:    "123",
	}
	err := card.Validate()
	s.Require().Nil(err)
	s.Require().Equal("1234 5678 1234 5678", card.Number)
}

func (s *CardSuite) TestInvalidName() {
	card := Card{
		Name:   "foo bar foo",
		Number: "12345678 1234  5678    ",
		CVV:    "123",
	}
	err := card.Validate()
	s.Require().ErrorIs(err, ErrInvalidName)
}

func (s *CardSuite) TestInvalidNumber() {
	card := Card{
		Name:   "foo bar",
		Number: "12345678 1234  5678    123",
		CVV:    "123",
	}
	err := card.Validate()
	s.Require().ErrorIs(err, ErrInvalidNumber)
}

func (s *CardSuite) TestInvalidCVV() {
	card := Card{
		Name:   "foo bar",
		Number: "12345678 1234  5678    ",
		CVV:    "1234",
	}
	err := card.Validate()
	s.Require().ErrorIs(err, ErrInvalidCVV)
}

func (s *CardSuite) TestEmptyCard() {
	card := Card{
		Name:   "",
		Number: "",
		CVV:    "",
	}
	s.Require().Equal(card.HasData(), false)
}

func (s *CardSuite) TestNonEmptyCard() {
	card := Card{
		Name:   "1",
		Number: "",
		CVV:    "",
	}
	s.Require().Equal(card.HasData(), true)
}
