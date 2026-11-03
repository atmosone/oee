package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atmosone/oee/internal/apps/cms/adapters/openapi/schema"
)

type AssetsCounterStore struct {
	db *sql.DB
}

// AssetsCounter implements [queries.AssetsCounterStore].
func (s *AssetsCounterStore) AssetsCounter(ctx context.Context) (schema.AssetsCounter, error) {
	const op = "apps.cms.adapters.postgres.AssetsCounterStore.AssetsCounter"
	var res schema.AssetsCounter
	if err := s.db.QueryRowContext(ctx,
		`SELECT 
			(SELECT COUNT(*) FROM equipment) AS total_equipment,
			(SELECT COUNT(*) FROM equipment_lines) AS total_equipment_lines,
			(SELECT COUNT(*) FROM specifications) AS total_specifications;`,
	).Scan(&res.TotalEquipment,
		&res.TotalEquipmentLines,
		&res.TotalSpecifications,
	); err != nil {
		return schema.AssetsCounter{}, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func NewAssetsCounterStore(db *sql.DB) *AssetsCounterStore {
	return &AssetsCounterStore{db: db}
}
