package postgres

import (
	"context"
	"database/sql"

	"github.com/atmosone/oee/internal/domain/workshifts"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Workshift(context.Context, workshifts.ID) (*workshifts.Workshift, error) {
	panic("unimplemented")
}

func (r *Repository) ActiveWorkshift(context.Context) (bool, error) {
	panic("unimplemented")
}

func (r *Repository) Save(context.Context, *workshifts.Workshift) error {
	panic("unimplemented")
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
