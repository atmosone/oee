package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

type CreateSpecificationRepository struct{ db *sql.DB }

// SaveSpecification implements [commands.CreateSpecificationRepository].
func (r *CreateSpecificationRepository) SaveSpecification(ctx context.Context, specification *domain.Specification) error {
	const op = "apps.cms.adapters.postgres.CreateSpecificationRepository.SaveSpecification"
	for _, e := range specification.Events() {
		if event, ok := e.(domain.SpecificationCreatedEvent); ok {
			tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			defer tx.Rollback()
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO specifications (id, code, description) VALUES ($1, $2, $3);`,
				event.SpecificationID,
				event.SpecificationCode,
				event.SpecificationDescription,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return nil
		}
		continue
	}
	return nil
}

func NewCreateSpecificationRepository(db *sql.DB) *CreateSpecificationRepository {
	return &CreateSpecificationRepository{db: db}
}
