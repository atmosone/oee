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
	//go:embed sql/read/workshift_query.sql
	workshiftQuery string
	//go:embed sql/write/start_workshift_query.sql
	startWorkshiftQuery string
	//go:embed sql/write/finish_workshift_query.sql
	finishWorkshiftQuery string
)

type WorkshiftRepository struct{ db *sql.DB }

func NewWorkshiftRepository(db *sql.DB) *WorkshiftRepository {
	return &WorkshiftRepository{db: db}
}

func (r *WorkshiftRepository) Workshift(
	ctx context.Context, workshiftID domain.WorkshiftID,
) (*domain.Workshift, error) {
	const op = "adapters.postgres.WorkshiftRepository.Workshift"
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var snapshot domain.WorkshiftSnapshot
	if err := tx.QueryRowContext(ctx, workshiftQuery, workshiftID).Scan(
		&snapshot.ID,
		&snapshot.Name,
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
	return domain.HydrateWorkshift(snapshot), nil
}

func (r *WorkshiftRepository) SaveWorkshift(
	ctx context.Context, workshift *domain.Workshift,
) error {
	const op = "adapters.postgres.WorkshiftRepository.SaveWorkshift"
	if events := workshift.Events(); events != nil {
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
			case domain.WorkshiftStartedEvent:
				if _, err := tx.ExecContext(ctx,
					startWorkshiftQuery,
					event.WorkshiftID,
					event.WorkshiftName,
					event.WorkshiftStatus,
					event.WorkshiftStartedAt,
					time.Time{}, // Костыльно, но без NULL
				); err != nil {
					return fmt.Errorf("%s: %w", op, err)
				}
			case domain.WorkshiftFinishedEvent:
				if _, err := tx.ExecContext(ctx,
					finishWorkshiftQuery,
					event.WorkshiftStatus,
					event.WorkshiftDuration,
					event.WorkshiftFinishedAt,
					event.WorkshiftID,
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
