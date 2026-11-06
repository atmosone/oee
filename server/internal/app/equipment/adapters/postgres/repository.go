package postgres

import (
	"context"
	"database/sql"

	"github.com/atmosone/oee/internal/domain/equipment"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Save(context.Context, *equipment.Equipment) error {
	panic("unimplemented")
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
