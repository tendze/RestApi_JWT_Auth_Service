package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"jwt-auth-service/internal/storage"
)

type Storage struct {
	DB *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "strg.postgresql.New"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    login VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(100) NOT NULL
)`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	strg := &Storage{DB: db}
	fmt.Println(strg.SaveUser("nnnn", "123123"))
	return strg, nil
}

func (s *Storage) SaveUser(login, password string) (int64, error) {
	const op = "storage.postgresql.SaveUser"
	var lastId int64
	err := s.DB.QueryRowContext(
		context.Background(),
		`INSERT INTO users(login, password) VALUES($1, $2) RETURNING id`,
		login, password,
	).Scan(&lastId)
	if err != nil {
		pgError, ok := err.(*pgconn.PgError)
		if ok && pgError.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return lastId, nil
}
