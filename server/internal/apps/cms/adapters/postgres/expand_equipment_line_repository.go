package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain"
)

type ExpandEquipmentLineRepository struct{ db *sql.DB }

// EquipmentLine implements [commands.ExpandEquipmentLineRepository].
func (r *ExpandEquipmentLineRepository) EquipmentLine(ctx context.Context, id domain.EquipmentLineID) (*domain.EquipmentLine, error) {
	const op = "apps.sms.adapters.postgres.ExpandEquipmentLineRepository.EquipmentLine"
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
		`SELECT id, name FROM equipment_lines WHERE id = $1;`,
		id,
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
		`SELECT equipment_id FROM equipment_line_items WHERE equipment_line_id = $1;`,
		id,
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

// SaveEquipmentLine implements [commands.ExpandEquipmentLineRepository].
func (r *ExpandEquipmentLineRepository) SaveEquipmentLine(ctx context.Context, el *domain.EquipmentLine) error {
	const op = "apps.cms.adapters.postgres.ExpandEquipmentLineRepository.SaveEquipmentLine"
	for _, e := range el.Events() {
		if event, ok := e.(domain.EquipmentLineExpandedEvent); ok {
			tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			if err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			defer tx.Rollback()
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO equipment_line_items (equipment_line_id, equipment_id) VALUES ($1, $2);`,
				event.EquipmentLineID,
				event.EquipmentID,
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

func NewExpandEquipmentLineRepository(db *sql.DB) *ExpandEquipmentLineRepository {
	return &ExpandEquipmentLineRepository{db: db}
}
