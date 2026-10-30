package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

var (
	//go:embed sql/write/create_equipment_query.sql
	createEquipmentQuery string
)

type EquipmentRepository struct{ db *sql.DB }

func NewEquipmentRepository(db *sql.DB) *EquipmentRepository {
	return &EquipmentRepository{db: db}
}

func (r *EquipmentRepository) SaveEquipment(ctx context.Context, equipment *domain.Equipment) error {
	const op = "adapters.postgres.EquipmentRepository.SaveEquipment"
	if events := equipment.Events(); events != nil {
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
			case domain.EquipmentCreatedEvent:
				if _, err := tx.ExecContext(ctx,
					createEquipmentQuery,
					event.EquipmentID,
					event.EquipmentName,
					event.EquipmentUnits,
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
