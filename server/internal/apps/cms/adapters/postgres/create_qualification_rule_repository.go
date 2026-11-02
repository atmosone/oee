package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

type CreateQualificationRuleRepository struct {
	db *sql.DB
}

// SaveQualificationRule implements [commands.CreateQualificationRuleRepository].
func (r *CreateQualificationRuleRepository) SaveQualificationRule(ctx context.Context, rule *domain.QualificationRule) error {
	const op = "apps.cms.adapters.postgres.CreateQualificationRuleRepository.SaveQualificationRule"
	for _, e := range rule.Events() {
		if event, ok := e.(domain.QualificationRuleCreatedEvent); ok {
			tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			defer tx.Rollback()
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO qualification_rules (equipment_id, specification_id) VALUES ($1, $2);`,
				event.EquipmentID,
				event.SpecificationID,
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

func NewCreateQualificationRuleRepository(db *sql.DB) *CreateQualificationRuleRepository {
	return &CreateQualificationRuleRepository{db: db}
}
