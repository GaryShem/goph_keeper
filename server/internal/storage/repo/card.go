package repo

import (
	"errors"
	"regexp"
)

var ErrInvalidName = errors.New("invalid name")
var ErrInvalidNumber = errors.New("invalid number")
var ErrInvalidCVV = errors.New("invalid CVV")

type Card struct {
	Name   string `json:"name"`
	Number string `json:"number"`
	CVV    string `json:"cvv"`
}

func (c *Card) Validate() error {
	nameRegexp, _ := regexp.Compile(`^[A-Za-z]+ [A-Za-z]+$`)
	if !nameRegexp.MatchString(c.Name) {
		return ErrInvalidName
	}
	numberRegexp, _ := regexp.Compile(`^(?:[0-9]\s*){16}$`)
	if !numberRegexp.MatchString(c.Number) {
		return ErrInvalidNumber
	}
	if len(c.Number) != len("1234 1234 1234 1234") {
		number := ""
		digits := 0
		for _, char := range c.Number {
			if char >= '0' && char <= '9' {
				number += string(char)
				digits++
				if digits%4 == 0 && digits < 15 {
					number += " "
				}
			}
		}
		c.Number = number
	}
	CVVRegexp, _ := regexp.Compile(`^\d{3}$`)
	if !CVVRegexp.MatchString(c.CVV) {
		return ErrInvalidCVV
	}
	return nil
}

func (c *Card) HasData() bool {
	return len(c.Name) > 0 || len(c.Number) > 0 || len(c.CVV) > 0
}
