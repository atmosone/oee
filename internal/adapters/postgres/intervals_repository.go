package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/intervals"
)

type IntervalReader struct{ db *sql.DB }

func NewIntervalReader(db *sql.DB) *IntervalReader {
	return &IntervalReader{db: db}
}

func (dr *IntervalReader) Equipment(ctx context.Context, id equipment.ID) (*intervals.Equipment, error) {
	const op = "adapters.postgres.intervals_repository.IntervalReader.Equipment"
	tx, err := dr.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	var equipment intervals.Equipment
	if err := tx.QueryRowContext(ctx,
		`SELECT id FROM equipment WHERE id = $1;`,
		id,
	).Scan(&equipment.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, intervals.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var active intervals.Interval
	if err := tx.QueryRowContext(ctx,
		`SELECT id, equipment_id, status, category, started_at 
		FROM intervals
		WHERE equipment_id = $1 
		AND status = $2;`,
		id, intervals.StatusActive,
	).Scan(&active.ID,
		&active.EquipmentID,
		&active.Status,
		&active.Category,
		&active.StartedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &equipment, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	equipment.Active = &active
	return &equipment, nil
}

type IntervalSaver struct{ db *sql.DB }

func NewIntervalSaver(db *sql.DB) *IntervalSaver {
	return &IntervalSaver{db: db}
}

func (ds *IntervalSaver) Save(ctx context.Context, equipment *intervals.Equipment) error {
	const op = "adapters.postgres.intervals_repository.IntervalSaver.Save"
	if len(equipment.Events) == 0 {
		return nil
	}
	tx, err := ds.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range equipment.Events {
		switch e := event.(type) {
		case intervals.StartedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO intervals (id, equipment_id, status, category, started_at, finished_at, duration) 
				VALUES ($1, $2, $3, $4, $5, $6, $7);`,
				e.ID,
				e.EquipmentID,
				e.Status,
				e.Category,
				e.StartedAt,
				time.Time{},
				0,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case intervals.FinishedEvent:
			if _, err := tx.ExecContext(ctx,
				`UPDATE intervals 
				SET status = $1, finished_at = $2, duration = $3 
				WHERE id = $4;`,
				e.Status, e.FinishedAt, e.Duration, e.ID,
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

type IntervalsRepository struct {
	*IntervalReader
	*IntervalSaver
}

func NewIntervalsRepository(
	reader *IntervalReader,
	saver *IntervalSaver,
) *IntervalsRepository {
	return &IntervalsRepository{
		IntervalReader: reader,
		IntervalSaver:  saver,
	}
}
