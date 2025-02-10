package sqlite

import (
	"database/sql"
	"errors"
	"url-shortener/internal/lib/e"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	q, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias url(alias);
	`)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	if _, err := q.Exec(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	const op = "storage.sqlite.SaveURL"

	q, err := s.db.Prepare(`INSERT INTO url(url, alias) VALUES (?, ?)`)
	if err != nil {
		return e.Wrap(op, err)
	}

	if _, err := q.Exec(urlToSave, alias); err != nil {
		sqliteErr, ok := err.(sqlite3.Error)
		isExists := ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique

		if isExists {
			return e.Wrap(op, storage.ErrURLExists)
		}

		return e.Wrap(op, err)
	}

	return nil
}

func (s *Storage) URL(alias string) (string, error) {
	const op = "storage.sqlite.URL"

	q, err := s.db.Prepare(`SELECT url FROM url WHERE alias=?`)
	if err != nil {
		return "", e.Wrap(op, err)
	}

	var resURL string

	err = q.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", e.Wrap(op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"

	q, err := s.db.Prepare(`delete from url where alias=?`)
	if err != nil {
		return e.Wrap(op, err)
	}

	if _, err := q.Exec(alias); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}
