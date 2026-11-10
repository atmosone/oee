package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
)

type DowntimeActiveChecker struct{ db *sql.DB }

func NewDowntimeActiveChecker(db *sql.DB) *DowntimeActiveChecker {
	return &DowntimeActiveChecker{db: db}
}

// ActiveDowntime implements [commands.StartDowntimeRepositoryPort].
func (dac *DowntimeActiveChecker) ActiveDowntime(ctx context.Context, equipmentID equipment.ID) (bool, error) {
	const op = "adapters.postgres.downtimes_repository.DowntimeActiveChecker.ActiveDowntime"
	var exists bool
	if err := dac.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM downtimes WHERE equipment_id = $1 AND status = $2);`,
		equipmentID, downtimes.StatusActive,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

type DowntimeReader struct{ db *sql.DB }

func NewDowntimeReader(db *sql.DB) *DowntimeReader {
	return &DowntimeReader{db: db}
}

func (dr *DowntimeReader) Downtime(ctx context.Context, id downtimes.ID) (*downtimes.Downtime, error) {
	const op = "adapters.postgres.downtimes_repository.DowntimeReader.Downtime"
	var downtime downtimes.Downtime
	if err := dr.db.QueryRowContext(ctx,
		`SELECT id, equipment_id, status, started_at, finished_at, duration FROM downtimes WHERE id = $1;`,
		id,
	).Scan(&downtime.ID,
		&downtime.EquipmentID,
		&downtime.Status,
		&downtime.StartedAt,
		&downtime.FinishedAt,
		&downtime.Duration,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, downtimes.ErrDowntimeNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &downtime, nil
}

type DowntimeSaver struct{ db *sql.DB }

func NewDowntimeSaver(db *sql.DB) *DowntimeSaver {
	return &DowntimeSaver{db: db}
}

func (ds *DowntimeSaver) Save(ctx context.Context, downtime *downtimes.Downtime) error {
	const op = "adapters.postgres.downtimes_repository.DowntimeSaver.Save"
	if len(downtime.Events) == 0 {
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
	for _, event := range downtime.Events {
		switch e := event.(type) {
		case downtimes.StartedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO downtimes (id, equipment_id, status, started_at, finished_at, duration) VALUES ($1, $2, $3, $4, $5, $6);`,
				e.ID, e.EquipmentID, e.Status, e.StartedAt, downtime.FinishedAt, downtime.Duration,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case downtimes.FinishedEvent:
			if _, err := tx.ExecContext(ctx,
				`UPDATE downtimes SET status = $1, finished_at = $2, duration = $3 WHERE id = $4;`,
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

type DowntimesRepository struct {
	*DowntimeActiveChecker
	*DowntimeReader
	*DowntimeSaver
}

func NewDowntimesRepository(
	activeChecker *DowntimeActiveChecker,
	reader *DowntimeReader,
	saver *DowntimeSaver,
) *DowntimesRepository {
	return &DowntimesRepository{
		DowntimeActiveChecker: activeChecker,
		DowntimeReader:        reader,
		DowntimeSaver:         saver,
	}
}
