package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

type SpecificationsListStore struct {
	db *sql.DB
}

// SpecificationsList implements [queries.SpecificationsListStore].
func (s *SpecificationsListStore) SpecificationsList(ctx context.Context) (schema.SpecificationsList, error) {
	const op = "apps.cms.adapters.postgres.SpecificationsListStore.SpecificationsList"
	rows, err := s.db.QueryContext(ctx,
		`SELECT s.id, s.code, s.description, COUNT(qr.equipment_id) AS total_qualification_rules
		FROM specifications s
		LEFT JOIN qualification_rules qr ON s.id = qr.specification_id
		GROUP BY s.id;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	res := make(schema.SpecificationsList, 0)
	for rows.Next() {
		var item schema.SpecificationsListItem
		if err := rows.Scan(
			&item.Id, &item.Code, &item.Description, &item.TotalQualificationRules,
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

func NewSpecificationsListStore(db *sql.DB) *SpecificationsListStore {
	return &SpecificationsListStore{db: db}
}
