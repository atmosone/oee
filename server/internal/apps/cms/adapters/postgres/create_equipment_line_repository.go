package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

type CreateEquipmentLineRepository struct {
	db *sql.DB
}

// SaveEquipmentLine implements [commands.CreateEquipmentLineRepository].
func (r *CreateEquipmentLineRepository) SaveEquipmentLine(ctx context.Context, equipmentLine *domain.EquipmentLine) error {
	const op = "apps.cms.adapters.postgres.CreateEquipmentLineRepository.SaveEquipmentLine"
	for _, e := range equipmentLine.Events() {
		if event, ok := e.(domain.EquipmentLineCreatedEvent); ok {
			tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			defer tx.Rollback()
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO equipment_lines (id, name) VALUES ($1, $2);`,
				event.EquipmentLineID,
				event.EquipmentLineName,
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

func NewCreateEquipmentLineRepository(db *sql.DB) *CreateEquipmentLineRepository {
	return &CreateEquipmentLineRepository{db: db}
}
