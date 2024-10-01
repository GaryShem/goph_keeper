package repo

import "errors"

var ErrInvalidDataType = errors.New("invalid data type")
var ErrHasMultipleData = errors.New("data has multiple types")
var ErrNoUser = errors.New("user missing")
var ErrNoName = errors.New("name missing")
var ErrNoDataToWrite = errors.New("no data to write")

type RepoData struct {
	User   string
	Name   string
	Type   DataType
	Text   string
	Card   Card
	Binary []byte
}

func (r *RepoData) Validate(hasValue bool) error {
	if r.User == "" {
		return ErrNoUser
	}
	if r.Name == "" {
		return ErrNoName
	}
	if !r.Type.IsValid() && hasValue {
		return ErrInvalidDataType
	}
	dataCount := 0
	if len(r.Text) > 0 {
		dataCount++
	}
	if len(r.Binary) > 0 {
		dataCount++
	}
	if r.Card.HasData() {
		dataCount++
	}
	if hasValue && dataCount == 0 {
		return ErrNoDataToWrite
	}
	if dataCount > 1 {
		return ErrHasMultipleData
	}
	switch r.Type {
	case TextType:
		if len(r.Text) == 0 && hasValue {
			return ErrNoDataToWrite
		}
	case BinaryType:
		if len(r.Binary) == 0 && hasValue {
			return ErrNoDataToWrite
		}
	case CardType:
		if err := r.Card.Validate(); err != nil && hasValue {
			return err
		}
	default:
		if hasValue {
			return ErrInvalidDataType
		}
	}

	return nil
}
