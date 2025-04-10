package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.NewStorage"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls (
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("не удалось подготовить оператор: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить оператор: %w", err)
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_alias ON urls (alias);`)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать индекс: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {

	const op = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare("INSERT INTO urls (alias, url) VALUES (?, ?);")
	if err != nil {
		return 0, fmt.Errorf("не удалось подготовить оператор: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			fmt.Printf("не удалось закрыть оператор: %v", err)
		}

	}(stmt)

	res, err := stmt.Exec(alias, urlToSave)
	if err != nil {

		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: дублирование уникального ключа: %w", op, storage.ErrURLAlreadyExists)
		}
		return 0, fmt.Errorf("не удалось выполнить оператор: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить последний вставленный id: %w", err)
	}
	return id, nil
}

func (s *Storage) GetURL(alies string) (string, error) {

	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?;")

	if err != nil {
		return "", fmt.Errorf("не удалось подготовить оператор: %w", err)
	}

	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			fmt.Printf("не удалось закрыть оператор: %v", err)
		}
	}(stmt)

	var url string
	err = stmt.QueryRow(alies).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return "", fmt.Errorf("не удалось выполнить оператор: %w", err)
	}
	return url, nil
}

func (s *Storage) DeleteUrl(alias string) error {

	const op = "storage.sqlite.DeleteUrl"

	stmt, err := s.db.Prepare("DELETE FROM urls WHERE alias = ?;")

	if err != nil {
		return fmt.Errorf("не удалось подготовить оператор: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			fmt.Printf("не удалось закрыть оператор: %v", err)
		}
	}(stmt)

	_, err = stmt.Exec(alias)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, storage.ErrURLNotFound)
		}
		return fmt.Errorf("не удалось выполнить оператор: %w", err)
	}
	return nil
}
