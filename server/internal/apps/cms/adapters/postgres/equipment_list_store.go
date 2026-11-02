package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

type EquipmentListStore struct {
	db *sql.DB
}

// EquipmentList implements [queries.EquipmentListStore].
func (s *EquipmentListStore) EquipmentList(ctx context.Context) (schema.EquipmentList, error) {
	const op = "apps.cms.adapters.postgres.EquipmentListStore.EquipmentList"
	rows, err := s.db.QueryContext(ctx,
		`SELECT e.id, e.name, e.units, COUNT(qr.specification_id) AS total_qualification_rules
		FROM equipment e
		LEFT JOIN qualification_rules qr ON e.id = qr.equipment_id
		GROUP BY e.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	res := make(schema.EquipmentList, 0)
	for rows.Next() {
		var item schema.EquipmentListItem
		if err := rows.Scan(
			&item.Id, &item.Name, &item.Units, &item.TotalQualificationRules,
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

func NewEquipmentListStore(db *sql.DB) *EquipmentListStore {
	return &EquipmentListStore{db: db}
}
