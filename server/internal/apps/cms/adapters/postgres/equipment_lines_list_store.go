package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

type EquipmentLinesListStore struct {
	db *sql.DB
}

// EquipmentLinesList implements [queries.EquipmentLinesListStore].
func (s *EquipmentLinesListStore) EquipmentLinesList(ctx context.Context) (schema.EquipmentLinesList, error) {
	const op = "apps.cms.adapters.postgres.EquipmentLinesListStore.EquipmentLinesList"
	rows, err := s.db.QueryContext(ctx,
		`SELECT el.id, el.name, COUNT(eli.equipment_id) AS total_equipment
		FROM equipment_lines el
		LEFT JOIN equipment_line_items eli ON el.id = eli.equipment_line_id
		GROUP BY el.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	res := make(schema.EquipmentLinesList, 0)
	for rows.Next() {
		var item schema.EquipmentLinesListItem
		if err := rows.Scan(
			&item.Id, &item.Name, &item.TotalEquipment,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func NewEquipmentLinesListStore(db *sql.DB) *EquipmentLinesListStore {
	return &EquipmentLinesListStore{db: db}
}
