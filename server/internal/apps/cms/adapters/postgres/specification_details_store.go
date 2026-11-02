package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
	"github.com/atmosone/oee/internal/apps/cms/queries"
)

type SpecificationDetailsStore struct {
	db *sql.DB
}

// SpecificationDetails implements [queries.SpecificationDetailsStore].
func (s *SpecificationDetailsStore) SpecificationDetails(ctx context.Context, params schema.SpecificationDetailsParams) (schema.SpecificationDetails, error) {
	const op = "apps.cms.adapters.postgres.SpecificationDetailsStore.SpecificationDetails"
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return schema.SpecificationDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var res schema.SpecificationDetails
	if err := tx.QueryRowContext(ctx,
		`SELECT code, description FROM specifications WHERE id = $1;`,
		params.Id,
	).Scan(&res.Code, &res.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return schema.SpecificationDetails{}, queries.ErrSpecificationDetailsNotFound
		}
		return schema.SpecificationDetails{}, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := tx.QueryContext(ctx,
		`SELECT e.id, e.name FROM equipment e
		JOIN qualification_rules qr ON e.id = qr.equipment_id
		WHERE qr.specification_id = $1;`,
		params.Id,
	)
	if err != nil {
		return schema.SpecificationDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var item schema.SpecificationDetailsEquipment
		if err := rows.Scan(
			&item.Id, &item.Name,
		); err != nil {
			return schema.SpecificationDetails{}, fmt.Errorf("%s: %w", op, err)
		}
		res.EquipmentList = append(res.EquipmentList, item)
	}
	if err := rows.Err(); err != nil {
		return schema.SpecificationDetails{}, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func NewSpecificationDetailsStore(db *sql.DB) *SpecificationDetailsStore {
	return &SpecificationDetailsStore{db: db}
}
