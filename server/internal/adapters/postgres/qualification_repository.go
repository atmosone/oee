package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

var (
	//go:embed sql/write/create_qualification_query.sql
	createQualificationQuery string
)

type QualificationRepository struct{ db *sql.DB }

func NewQualificationRepository(db *sql.DB) *QualificationRepository {
	return &QualificationRepository{db: db}
}

// SaveQualification implements [app.Repository].
func (r *QualificationRepository) SaveQualification(ctx context.Context, qualification *domain.Qualification) error {
	const op = "adapters.postgres.QualificationRepository.SaveQualification"
	if events := qualification.Events(); events != nil {
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
			case domain.QualificationCreatedEvent:
				if _, err := tx.ExecContext(ctx,
					createQualificationQuery,
					event.DowntimeID,
					event.SpecificationID,
					event.QualificationComment,
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
