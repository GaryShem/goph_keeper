package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"goph_keeper/goph_server/internal/config"
	"goph_keeper/goph_server/internal/logging"
	"goph_keeper/goph_server/internal/storage/repo"
)

type PostgreSQLRepo struct {
	db   *sql.DB
	lock *sync.RWMutex
}

func (p *PostgreSQLRepo) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 300*time.Second)
}

func NewPostgreSQLRepo(config config.ServerConfig) (*PostgreSQLRepo, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DbHost, config.DbPort, config.DbUser, config.DbPass, config.DbName)
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		logging.Log().Info("cannot connect to db", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err = tx.Exec(createUserTableSQL); err != nil {
		return nil, err
	}
	if _, err = tx.Exec(createDataTableQSL); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &PostgreSQLRepo{db: db, lock: &sync.RWMutex{}}, nil
}

var _ repo.Repo = (*PostgreSQLRepo)(nil)
