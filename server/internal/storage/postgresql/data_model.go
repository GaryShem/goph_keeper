package postgresql

import (
	"database/sql"
	"encoding/json"
	"errors"

	"goph_keeper/goph_server/internal/storage/repo"
)

var ErrDifferentDataAlreadyExists = errors.New("data with that name but different type already exists")
var ErrDataNotFound = errors.New("data not found")

type DbData struct {
	User     string
	Name     string
	DataType repo.DataType
	Text     sql.NullString
	CardJson sql.NullString
	Binary   []byte
}

func FromRepoData(rData repo.RepoData) (*DbData, error) {
	textData := sql.NullString{
		String: rData.Text,
		Valid:  rData.Text != "",
	}
	cardData := sql.NullString{}
	if !rData.Card.HasData() {
		cardData.Valid = false
		cardData.String = ""
	} else {
		cardJson, err := json.Marshal(rData.Card)
		if err != nil {
			return nil, err
		}
		cardData.String = string(cardJson)
		cardData.Valid = true
	}

	result := &DbData{
		User:     rData.User,
		Name:     rData.Name,
		DataType: rData.Type,
		Text:     textData,
		CardJson: cardData,
		Binary:   rData.Binary,
	}
	return result, nil

}

func (d *DbData) ToRepoData() (*repo.RepoData, error) {
	text := ""
	if d.Text.Valid {
		text = d.Text.String
	}
	card := repo.Card{}
	if d.CardJson.Valid {
		err := json.Unmarshal([]byte(d.CardJson.String), &card)
		if err != nil {
			return nil, err
		}
	}
	result := repo.RepoData{
		User:   d.User,
		Name:   d.Name,
		Type:   d.DataType,
		Text:   text,
		Card:   card,
		Binary: d.Binary,
	}
	return &result, nil
}

func (d *DbData) ValidateRead() error {
	rData, err := d.ToRepoData()
	if err != nil {
		return err
	}
	return rData.Validate(false)
}

func (d *DbData) ValidateWrite() error {
	rData, err := d.ToRepoData()
	if err != nil {
		return err
	}
	return rData.Validate(true)
}
