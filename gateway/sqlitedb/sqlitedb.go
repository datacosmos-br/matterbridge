// Package sqlitedb opens the sqlite databases kept by the gateway stores.
package sqlitedb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // database/sql driver "sqlite"
)

// Open opens (creating if needed) the sqlite database at path in WAL mode
// and applies schema.
func Open(path, schema string) (*sql.DB, error) {
	err := os.MkdirAll(filepath.Dir(path), 0o750)
	if err != nil {
		return nil, fmt.Errorf("create directory of %s: %w", path, err)
	}

	dsn := "file:" + path + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)"

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	_, err = db.ExecContext(context.Background(), schema)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("create schema in %s: %w", path, err), db.Close())
	}

	return db, nil
}

// Prepare returns a prepared statement owned by the caller.
func Prepare(db *sql.DB, query string) (*sql.Stmt, error) {
	stmt, err := db.PrepareContext(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("prepare %q: %w", query, err)
	}

	return stmt, nil
}

// CloseAll closes the statements that were prepared and then the database.
func CloseAll(db *sql.DB, stmts ...*sql.Stmt) error {
	var errs []error

	for _, stmt := range stmts {
		if stmt != nil {
			errs = append(errs, stmt.Close())
		}
	}

	errs = append(errs, db.Close())

	return errors.Join(errs...)
}
