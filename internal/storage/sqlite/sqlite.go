package sqlite

import (
	"cmd/service/internal/storage"
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	op := "storage.sqlite.SaveUrl"
	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	//TODO : error is not clear, need refactoring
	if err != nil {
		if sqliteErr, ok := err.(*sqlite.Error); ok && sqliteErr.Code() == 19 {
			return 0, fmt.Errorf("%s:%w", op, storage.ErrUrlExist)
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	//TODO : LastInsertId is platform dependant, doesnt work in MSQL
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s failed to get last inserted id :%w", op, err)
	}
	return id, nil

}

func (s *Storage) GetUrl(alias string) (string, error) {
	op := "storage.sqlite.GetUrl"
	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	var retURL string
	err = stmt.QueryRow(alias).Scan(&retURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s failed to execute statement :%w", op, err)
	}
	return retURL, nil
}
