package postgres

import (
	"context"
	"database/sql"

	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
)

type Repository struct {
	db *sql.DB
}

// Downtime implements [commands.FinishRepositoryPort].
func (r *Repository) Downtime(context.Context, downtimes.ID) (*downtimes.Downtime, error) {
	panic("unimplemented")
}

// ActiveDowntime implements [commands.StartRepositoryPort].
func (r *Repository) ActiveDowntime(context.Context, equipment.ID) (bool, error) {
	panic("unimplemented")
}

// Save implements [commands.StartRepositoryPort].
func (r *Repository) Save(context.Context, *downtimes.Downtime) error {
	panic("unimplemented")
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
