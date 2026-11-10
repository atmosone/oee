package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/domain/equipment"
)

type EquipmentSaver struct{ db *sql.DB }

func NewEquipmentSaver(db *sql.DB) *EquipmentSaver { return &EquipmentSaver{db: db} }

func (es *EquipmentSaver) Save(ctx context.Context, eq *equipment.Equipment) error {
	const op = "adapters.postgres.equipment_repository.EquipmentSaver.Save"
	tx, err := es.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range eq.Events {
		switch e := event.(type) {
		case equipment.CreatedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO equipment (id, name, category) VALUES ($1, $2, $3);`,
				e.ID, e.Name, e.Category,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		default:
			return nil
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type EquipmentRepository struct {
	*EquipmentSaver
}

func NewEquipmentRepository(
	saver *EquipmentSaver,
) *EquipmentRepository {
	return &EquipmentRepository{
		EquipmentSaver: saver,
	}
}
