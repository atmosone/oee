package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
	"github.com/atmosone/oee/internal/apps/cms/queries"
)

type EquipmentLineDetailsStore struct {
	db *sql.DB
}

// EquipmentLineDetails implements [queries.EquipmentLineDetailsStore].
func (s *EquipmentLineDetailsStore) EquipmentLineDetails(ctx context.Context, params schema.EquipmentLineDetailsParams) (schema.EquipmentLineDetails, error) {
	const op = "apps.cms.adapters.postgres.EquipmentLineDetailsStore.EquipmentLineDetails"
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return schema.EquipmentLineDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var res schema.EquipmentLineDetails
	if err := tx.QueryRowContext(ctx,
		`SELECT name FROM equipment_lines WHERE id = $1;`,
		params.Id,
	).Scan(&res.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schema.EquipmentLineDetails{}, queries.ErrEquipmentLineDetailsNotFound
		}
		return schema.EquipmentLineDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.QueryContext(ctx,
		`SELECT e.id, e.name FROM equipment e
		JOIN equipment_line_items eli ON e.id = eli.equipment_id
		WHERE eli.equipment_line_id = $1;`,
		params.Id,
	)
	if err != nil {
		return schema.EquipmentLineDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var item schema.EquipmentLineDetailsEquipment
		if err := rows.Scan(
			&item.Id, &item.Name,
		); err != nil {
			return schema.EquipmentLineDetails{}, fmt.Errorf("%s: %w", op, err)
		}
		res.EquipmentList = append(res.EquipmentList, item)
	}
	if err := rows.Err(); err != nil {
		return schema.EquipmentLineDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func NewEquipmentLineDetailsStore(db *sql.DB) *EquipmentLineDetailsStore {
	return &EquipmentLineDetailsStore{db: db}
}
