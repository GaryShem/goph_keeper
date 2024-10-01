package postgresql

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"goph_keeper/goph_server/internal/storage/repo"
)

func (p *PostgreSQLRepo) DownloadStorage(user string) ([]repo.RepoData, error) {
	result := make([]repo.RepoData, 0)
	ctx, cancel := p.getContext()
	defer cancel()
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	rows, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("dataname", "data_type", "text_data", "card_data", "blobsky_data").
		From("data").
		Where(sq.Eq{"username": user}).
		RunWith(tx).Query()
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		data := DbData{User: user}
		if err = rows.Scan(&data.Name, &data.DataType, &data.Text, &data.CardJson, &data.Binary); err != nil {
			return nil, err
		}
		repoData, rowErr := data.ToRepoData()
		if rowErr != nil {
			return nil, rowErr
		}
		result = append(result, *repoData)
	}
	return result, nil
}
