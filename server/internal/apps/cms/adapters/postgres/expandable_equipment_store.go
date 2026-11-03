package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

type ExpandableEquipmentStore struct {
	db *sql.DB
}

// ExpandableEquipment implements [queries.ExpandableEquipmentStore].
func (s *ExpandableEquipmentStore) ExpandableEquipment(ctx context.Context, params schema.ExpandableEquipmentParams) (schema.ExpandableEquipmentList, error) {
	const op = "apps.cms.adapters.postgres.EquipmentListStore.EquipmentList"
	rows, err := s.db.QueryContext(ctx,
		`SELECT e.id, e.name 
		FROM equipment e 
		WHERE NOT EXISTS (
    		SELECT 1 
    		FROM equipment_line_items eli 
    		WHERE eli.equipment_id = e.id 
		);`,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	res := make(schema.ExpandableEquipmentList, 0)
	for rows.Next() {
		var item schema.ExpandableEquipment
		if err := rows.Scan(&item.Id, &item.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil

}

func NewExpandableEquipmentStore(db *sql.DB) *ExpandableEquipmentStore {
	return &ExpandableEquipmentStore{db: db}
}
