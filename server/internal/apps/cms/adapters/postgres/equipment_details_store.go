package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
	"github.com/atmosone/oee/internal/apps/cms/queries"
)

type EquipmentDetailsStore struct {
	db *sql.DB
}

// EquipmentDetails implements [queries.EquipmentDetailsStore].
func (s *EquipmentDetailsStore) EquipmentDetails(ctx context.Context, params schema.EquipmentDetailsParams) (schema.EquipmentDetails, error) {
	const op = "apps.cms.adapters.postgres.EquipmentDetailsStore.EquipmentDetails"
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return schema.EquipmentDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var res schema.EquipmentDetails
	if err := tx.QueryRowContext(ctx,
		`SELECT name, units FROM equipment WHERE id = $1;`,
		params.Id,
	).Scan(
		&res.Name, &res.Units,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schema.EquipmentDetails{}, queries.ErrEquipmentDetailsNotFound
		}
		return schema.EquipmentDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	var line schema.EquipmentDetailsLine
	if err := tx.QueryRowContext(ctx,
		`SELECT el.id, el.name FROM equipment_lines el 
		JOIN equipment_line_items eli ON el.id = eli.equipment_line_id
		WHERE eli.equipment_id = $1;`,
		params.Id,
	).Scan(
		&line.Id, &line.Name,
	); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return schema.EquipmentDetails{}, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		res.Line = &line
	}

	rows, err := tx.QueryContext(ctx,
		`SELECT s.id, s.code, s.description FROM specifications s
		JOIN qualification_rules qr ON s.id = qr.specification_id
		WHERE qr.equipment_id = $1;`,
		params.Id,
	)
	if err != nil {
		return schema.EquipmentDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var item schema.EquipmentDetailsSpecification
		if err := rows.Scan(
			&item.Id, &item.Code, &item.Description,
		); err != nil {
			return schema.EquipmentDetails{}, fmt.Errorf("%s: %w", op, err)
		}
		res.SpecificationsList = append(res.SpecificationsList, item)
	}
	if err := rows.Err(); err != nil {
		return schema.EquipmentDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func NewEquipmentDetailsStore(db *sql.DB) *EquipmentDetailsStore {
	return &EquipmentDetailsStore{db: db}
}
