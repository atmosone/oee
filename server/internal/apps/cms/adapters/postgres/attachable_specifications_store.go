package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

type AttachableSpecificationsStore struct {
	db *sql.DB
}

// AttachableSpecifications implements [queries.AttachableSpecificationsStore].
func (r *AttachableSpecificationsStore) AttachableSpecifications(
	ctx context.Context, params schema.AttachableSpecificationsParams,
) (schema.AttachableSpecificationsList, error) {
	const op = "apps.cms.adapters.postgres.AttachableSpecificationsStore.AttachableSpecifications"
	rows, err := r.db.QueryContext(ctx,
		`SELECT s.id, s.code, s.description 
		FROM specifications s 
		WHERE NOT EXISTS (
    		SELECT 1 
    		FROM qualification_rules qr 
    		WHERE qr.equipment_id = $1 
      		AND qr.specification_id = s.id
		);`,
		params.To,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	res := make(schema.AttachableSpecificationsList, 0)
	for rows.Next() {
		var item schema.AttachableSpecification
		if err := rows.Scan(
			&item.Id, &item.Code, &item.Description,
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

func NewAttachableSpecificationsStore(db *sql.DB) *AttachableSpecificationsStore {
	return &AttachableSpecificationsStore{db: db}
}
