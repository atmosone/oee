package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

var (
	//go:embed sql/write/create_qualification_rule_query.sql
	createQualificationRuleQuery string
)

type QualificationRuleRepository struct{ db *sql.DB }

func NewQualificationRuleRepository(db *sql.DB) *QualificationRuleRepository {
	return &QualificationRuleRepository{db: db}
}

// SaveQualificationRule implements [app.Repository].
func (r *QualificationRuleRepository) SaveQualificationRule(ctx context.Context, qualificationRule *domain.QualificationRule) error {
	const op = "adapters.postgres.QualificationRuleRepository.SaveQualificationRule"
	if events := qualificationRule.Events(); events != nil {
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
			case domain.QualificationRuleCreatedEvent:
				if _, err := tx.ExecContext(ctx,
					createQualificationRuleQuery,
					event.EquipmentID,
					event.SpecificationID,
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
