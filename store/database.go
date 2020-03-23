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

	_, err = db.Exec(`
		create table if not exists releases (
			provider text, project text, version text,
			primary key (provider, project, version)
		)
	`)
	if err != nil {
		return nil, err
	}

	store := SQLiteStore{db: db}
	return &store, nil
}

func (s *SQLiteStore) Exists(release model.Release) bool {
	row := s.db.QueryRow(`
		select 1 from releases where provider = ? and project = ? and version = ?
	`, release.Provider.GetName(), release.Project.String(), release.Name)

	var result int
	row.Scan(&result)
	return result != 0
}

func (s *SQLiteStore) Mark(release model.Release) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		insert into releases (provider, project, version) values (?, ?, ?)
	`, release.Provider.GetName(), release.Project.String(), release.Name)
	if err == nil {
		tx.Commit()
	}

	return err
}

func (s *SQLiteStore) Close() {
	s.db.Close()
}
