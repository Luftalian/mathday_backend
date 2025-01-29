package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// BeginTransaction starts a new database transaction

func (r *Repository) BeginTransaction(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
