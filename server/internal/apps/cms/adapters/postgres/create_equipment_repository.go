package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

type CreateEquipmentRepository struct {
	db *sql.DB
}

// SaveEquipment implements [commands.CreateEquipmentRepository].
func (r *CreateEquipmentRepository) SaveEquipment(ctx context.Context, equipment *domain.Equipment) error {
	const op = "apps.cms.adapters.postgres.CreateEquipmentRepository.SaveEquipment"
	for _, e := range equipment.Events() {
		if event, ok := e.(domain.EquipmentCreatedEvent); ok {
			tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			defer tx.Rollback()
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO equipment (id, name, units) VALUES ($1, $2, $3);`,
				event.EquipmentID,
				event.EquipmentName,
				event.EquipmentUnits,
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

func NewCreateEquipmentRepository(db *sql.DB) *CreateEquipmentRepository {
	return &CreateEquipmentRepository{db: db}
}
