package postgres

import (
	"context"
	"database/sql"

	"github.com/atmosone/oee/internal/domain/specifications"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Specification(context.Context, specifications.ID) (*specifications.Specification, error) {
	panic("unimplemented")
}

func (r *Repository) Save(context.Context, *specifications.Specification) error {
	panic("unimplemented")
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
