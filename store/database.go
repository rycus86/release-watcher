package store

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rycus86/release-watcher/model"
)

type SQLiteStore struct {
	db *sql.DB
}

func Initialize(path string) (model.Store, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("create table if not exists latest_version (name text primary key, version text)")
	if err != nil {
		return nil, err
	}

	store := SQLiteStore{db: db}
	return &store, nil
}

func (s *SQLiteStore) Get(key string) string {
	row := s.db.QueryRow("select version from latest_version where name = ?", key)
	if row != nil {
		var value string
		row.Scan(&value)
		return value
	}

	return ""
}

func (s *SQLiteStore) Set(key string, value string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("insert or replace into latest_version (name, version) values (?, ?)", key, value)
	if err == nil {
		tx.Commit()
	}

	return err
}

func (s *SQLiteStore) Close() {
	s.db.Close()
}
