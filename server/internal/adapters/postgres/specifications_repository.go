package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

var (
	//go:embed sql/write/create_specification_query.sql
	createSpecificationQuery string
)

type SpecificationsRepository struct{ db *sql.DB }

func NewSpecificationsRepository(db *sql.DB) *SpecificationsRepository {
	return &SpecificationsRepository{db: db}
}

// SaveSpecification implements [app.Repository].
func (r *SpecificationsRepository) SaveSpecification(ctx context.Context, specification *domain.Specification) error {
	const op = "adapters.postgres.SpecificationsRepository.SaveSpecification"
	if events := specification.Events(); events != nil {
		tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer tx.Rollback()
		for _, e := range events {
			switch event := e.(type) {
			case domain.SpecificationCreatedEvent:
				if _, err := tx.ExecContext(ctx,
					createSpecificationQuery,
					event.SpecificationID,
					event.SpecificationCode,
					event.SpecificationDescription,
				); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			}
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
