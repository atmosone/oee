package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

var (
	//go:embed sql/read/equipment_line_query.sql
	equipmentLineQuery string
	//go:embed sql/read/equipment_line_items_query.sql
	equipmentLineItemsQuery string
	//go:embed sql/write/create_equipment_line_query.sql
	createEquipmentLineQuery string
	//go:embed sql/write/expand_equipment_line_query.sql
	expandEquipmentLineQuery string
)

type EquipmentLineRepository struct{ db *sql.DB }

func NewEquipmentLineRepository(db *sql.DB) *EquipmentLineRepository {
	return &EquipmentLineRepository{db: db}
}

// EquipmentLine implements [app.Repository].
func (r *EquipmentLineRepository) EquipmentLine(
	ctx context.Context, equipmentLineID domain.EquipmentLineID,
) (*domain.EquipmentLine, error) {
	const op = "adapters.postgres.EquipmentLineRepository.EquipmentLine"
	var snapshot domain.EquipmentLineSnapshot
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	if err := tx.QueryRowContext(ctx,
		equipmentLineQuery,
		equipmentLineID,
	).Scan(
		&snapshot.ID,
		&snapshot.Name,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := tx.QueryContext(ctx,
		equipmentLineItemsQuery,
		equipmentLineID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		snapshot.Equipment = append(snapshot.Equipment, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return domain.HydrateEquipmentLine(snapshot), nil
}

func (r *EquipmentLineRepository) SaveEquipmentLine(
	ctx context.Context, equipmentLine *domain.EquipmentLine,
) error {
	const op = "adapters.postgres.EquipmentLineRepository.SaveEquipmentLine"
	if events := equipmentLine.Events(); events != nil {
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
			case domain.EquipmentLineCreatedEvent:
				if _, err := tx.ExecContext(ctx,
					createEquipmentLineQuery,
					event.EquipmentLineID,
					event.EquipmentLineName,
				); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			case domain.EquipmentLineExpandedEvent:
				if _, err := tx.ExecContext(ctx,
					expandEquipmentLineQuery,
					event.EquipmentLineID,
					event.EquipmentID,
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
