package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
)

var ErrUnableRegister = errors.New("unable to register user")
var ErrUserAlreadyExists = errors.New("user already exists")

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

func (p *PostgreSQLRepo) LoginUser(user, pass string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	query := sq.Select("username", "password").From("users").Where(sq.Eq{"username": user}).
		PlaceholderFormat(sq.Dollar)
	var username, password string
	err = query.RunWith(tx).QueryRow().Scan(&username, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		} else {
			return err
		}
	}
	if pass != password {
		return ErrInvalidPassword
	}
	return nil
}

func (p *PostgreSQLRepo) RegisterUser(user, pass string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	uniqueQuery := sq.Select("username").From("users").Where(sq.Eq{"username": user}).PlaceholderFormat(sq.Dollar)
	err = uniqueQuery.RunWith(tx).QueryRow().Scan()
	if err == nil {
		// we got a row with this username, return error
		return fmt.Errorf(`%w: %w`, ErrUnableRegister, ErrUserAlreadyExists)
	} else {
		if errors.Is(err, sql.ErrNoRows) {
			// user with this name does not exist, continue
		} else {
			// generic unexpected error, return raw
			return fmt.Errorf(`%w: %w`, ErrUnableRegister, err)
		}
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("users").
		Columns("username", "password").
		Values(user, pass)
	if _, err = query.RunWith(p.db).Query(); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf(`%w: %w`, ErrUnableRegister, err)
	}
	return nil
}
