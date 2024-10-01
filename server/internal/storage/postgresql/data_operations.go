package postgresql

import (
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"goph_keeper/goph_server/internal/storage/repo"
)

func (p *PostgreSQLRepo) GetData(repoData repo.RepoData) (*repo.RepoData, error) {
	data, err := FromRepoData(repoData)
	if err != nil {
		return nil, err
	}
	if err = data.ValidateRead(); err != nil {
		return nil, err
	}
	result := DbData{
		User: data.User,
		Name: data.Name,
	}
	ctx, cancel := p.getContext()
	defer cancel()
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return nil, err
	}
	err = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("data_type", "text_data", "card_data", "blobsky_data").
		From("data").
		Where(sq.Eq{"username": data.User, "dataname": data.Name}).
		RunWith(tx).QueryRow().Scan(&result.DataType, &result.Text, &result.CardJson, &result.Binary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDataNotFound
		}
		return nil, err
	}
	repoResult, err := result.ToRepoData()
	if err != nil {
		return nil, err
	}
	return repoResult, nil
}

func (p *PostgreSQLRepo) checkUniqueData(tx *sql.Tx, data DbData) error {
	var dataType repo.DataType

	err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("data_type").
		From("data").
		Where(sq.Eq{"username": data.User, "dataname": data.Name}).
		RunWith(tx).QueryRow().Scan(&dataType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// data with this name does not exist, all good
			return nil
		} else {
			// unexpected generic error, return as is
			return err
		}
	}
	if dataType != data.DataType {
		return ErrDifferentDataAlreadyExists
	}
	return nil
}

func (p *PostgreSQLRepo) SetData(repoData repo.RepoData) error {
	data, err := FromRepoData(repoData)
	if err != nil {
		return err
	}
	if err = data.ValidateWrite(); err != nil {
		return err
	}
	ctx, cancel := p.getContext()
	defer cancel()
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if err = p.checkUniqueData(tx, *data); err != nil {
		return err
	}

	_, err = sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert("data").
		Columns("username", "dataname", "data_type", "text_data", "card_data", "blobsky_data").
		Values(data.User, data.Name, data.DataType, data.Text, data.CardJson, data.Binary).
		Suffix("ON CONFLICT (username, dataname) DO UPDATE SET text_data = EXCLUDED.text_data, card_data = EXCLUDED.card_data, blobsky_data = EXCLUDED.blobsky_data").
		RunWith(tx).Query()
	if err != nil {
		return err
	}
	return tx.Commit()
}
