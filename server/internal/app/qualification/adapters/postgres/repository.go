package postgres

import (
	"context"
	"database/sql"

	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Downtime(ctx context.Context, id downtimes.ID) (qualification.Downtime, error) {
	panic("unimplemented")
}

func (r *Repository) Specification(ctx context.Context, id specifications.ID) (qualification.Specification, error) {
	panic("unimplemented")
}

func (r *Repository) Qualification(ctx context.Context, id equipment.ID) (*qualification.Qualification, error) {
	panic("unimplemented")
}

func (r *Repository) Save(ctx context.Context, qualification *qualification.Qualification) error {
	panic("unimplemented")
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
