package repository

import (
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"

	"iziBookTest/internal/config"
)

type PgStorage struct {
	usersDB     *sqlx.DB
	documentsDB *sqlx.DB
}

func connectionString(User, Pass, Host, Port, DBName string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%v/%s?sslmode=disable",
		User, Pass, Host, Port, DBName,
	)
}

func NewStorage(conf *config.DB) (*PgStorage, error) {
	userConfig, err := pgx.ParseConfig(
		connectionString(conf.User, conf.Password, conf.Host, conf.Port, conf.UsersDbName),
	)
	if err != nil {
		return nil, err
	}

	userDB := stdlib.OpenDB(*userConfig)
	err = userDB.Ping()
	if err != nil {
		return nil, err
	}

	documentConfig, err := pgx.ParseConfig(
		connectionString(conf.User, conf.Password, conf.Host, conf.Port, conf.DocumentsDbName),
	)
	if err != nil {
		return nil, err
	}

	documentDB := stdlib.OpenDB(*documentConfig)
	err = documentDB.Ping()
	if err != nil {
		return nil, err
	}

	return &PgStorage{
		usersDB:     sqlx.NewDb(userDB, "postgres"),
		documentsDB: sqlx.NewDb(documentDB, "postgres"),
	}, nil
}
