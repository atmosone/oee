package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/intervals"
	"github.com/atmosone/oee/internal/domain/production"
)

type ProductionReader struct {
	db *sql.DB
}

func NewProductionReader(db *sql.DB) *ProductionReader {
	return &ProductionReader{db: db}
}

func (pr *ProductionReader) Counter(ctx context.Context,
	equipmentID equipment.ID,
	category production.Category,
) (*production.Counter, error) {
	const op = "adapters.postgres.production_repository.ProductionReader.Counter"
	tx, err := pr.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var counter production.Counter
	if err := tx.QueryRowContext(ctx,
		`SELECT id, equipment_id, category FROM intervals
		WHERE equipment_id = $1 
		AND status = $2 
		AND category = $3;`,
		equipmentID, intervals.StatusActive, intervals.CategoryWorktime,
	).Scan(&counter.Interval.ID,
		&counter.Interval.EquipmentID,
		&counter.Interval.Category,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, production.ErrActiveIntervalNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	counter.Production.Category = category
	if err := tx.QueryRowContext(ctx,
		`SELECT amount FROM production WHERE interval_id = $1 AND category = $2;`,
		counter.Interval.ID, counter.Production.Category,
	).Scan(&counter.Production.Amount); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &counter, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &counter, nil
}

type ProductionSaver struct {
	db *sql.DB
}

func NewProductionSaver(db *sql.DB) *ProductionSaver {
	return &ProductionSaver{db: db}
}

func (ps *ProductionSaver) Save(ctx context.Context, counter *production.Counter) error {
	const op = "adapters.postgres.production_repository.ProductionSaver.Save"
	if len(counter.Events) == 0 {
		return nil
	}
	tx, err := ps.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range counter.Events {
		switch e := event.(type) {
		case production.IncreasedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO production (interval_id, category, amount) VALUES ($1, $2, $3) 
				ON CONFLICT (interval_id, category) DO UPDATE SET amount = EXCLUDED.amount;`,
				e.IntervalID, e.Category, e.Amount,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type ProductionRepository struct {
	*ProductionReader
	*ProductionSaver
}

func NewProductionRepository(
	reader *ProductionReader,
	saver *ProductionSaver,
) *ProductionRepository {
	return &ProductionRepository{
		ProductionReader: reader,
		ProductionSaver:  saver,
	}
}
