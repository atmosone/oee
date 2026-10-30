package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	"github.com/atmosone/oee/internal/domain"
)

var (
	//go:embed sql/read/pull_qualified_downtime_equipment_query.sql
	pullQualifiedDowntimeEquipmentQuery string
	//go:embed sql/read/downtime_query.sql
	downtimeQuery string
	//go:embed sql/write/start_downtime_query.sql
	startDowntimeQuery string
	//go:embed sql/write/finish_downtime_query.sql
	finishDowntimeQuery string
)

type DowntimeRepository struct{ db *sql.DB }

func NewDowntimeRepository(db *sql.DB) *DowntimeRepository {
	return &DowntimeRepository{db: db}
}

func (r *DowntimeRepository) Downtime(ctx context.Context, downtimeID domain.DowntimeID) (*domain.Downtime, error) {
	const op = "adapters.postgres.DowntimeRepository.Downtime"
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var snapshot domain.DowntimeSnapshot
	if err := tx.QueryRowContext(ctx, downtimeQuery, downtimeID).Scan(
		&snapshot.ID,
		&snapshot.EquipmentID,
		&snapshot.Status,
		&snapshot.Duration,
		&snapshot.StartedAt,
		&snapshot.FinishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return domain.HydrateDowntime(snapshot), nil
}

func (r *DowntimeRepository) SaveDowntime(ctx context.Context, downtime *domain.Downtime) error {
	const op = "adapters.postgres.DowntimeRepository.SaveDowntime"
	if events := downtime.Events(); events != nil {
		tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		defer tx.Rollback()
		for _, e := range events {
			switch event := e.(type) {
			case domain.DowntimeStartedEvent:
				if _, err := tx.ExecContext(ctx,
					startDowntimeQuery,
					event.DowntimeID,
					event.EquipmentID,
					event.DowntimeStatus,
					event.DowntimeStartedAt,
					time.Time{}, // Костыльно, но без NULL
				); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			case domain.DowntimeFinishedEvent:
				if _, err := tx.ExecContext(ctx,
					finishDowntimeQuery,
					event.DowntimeStatus,
					event.DowntimeDuration,
					event.DowntimeFinishedAt,
					event.DowntimeID,
				); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			}
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
