package postgres

import (
	"context"
	"database/sql"

	"github.com/atmosone/oee/internal/domain/lines"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Line(context.Context, lines.ID) (*lines.Line, error) {
	panic("unimplemented")
}

func (r *Repository) Save(context.Context, *lines.Line) error {
	panic("unimplemented")
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
