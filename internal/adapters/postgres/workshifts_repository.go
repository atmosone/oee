package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain/workshifts"
)

type WorkshiftActiveChecker struct {
	db *sql.DB
}

func NewWorkshiftActiveChecker(db *sql.DB) *WorkshiftActiveChecker {
	return &WorkshiftActiveChecker{db: db}
}

func (wac *WorkshiftActiveChecker) ActiveWorkshift(ctx context.Context) (bool, error) {
	const op = "adapters.postgres.workshifts_repository.WorkshiftActiveChecker.ActiveWorkshift"
	var exists bool
	if err := wac.db.QueryRowContext(ctx,
		`SELECT EXISTS (SELECT 1 FROM workshifts WHERE status = $1);`,
		workshifts.StatusActive,
	).Scan(&exists); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return exists, nil
}

type WorkshiftReader struct {
	db *sql.DB
}

func NewWorkshiftReader(db *sql.DB) *WorkshiftReader {
	return &WorkshiftReader{db: db}
}

func (wr *WorkshiftReader) Workshift(ctx context.Context, id workshifts.ID) (*workshifts.Workshift, error) {
	const op = "adapters.postgres.workshifts_repository.WorkshiftReader.Workshift"
	var workshift workshifts.Workshift
	if err := wr.db.QueryRowContext(ctx,
		`SELECT id, name, status, started_at, finished_at, duration FROM workshfts WHERE id = $1;`,
		id,
	).Scan(&workshift.ID,
		&workshift.Name,
		&workshift.Status,
		&workshift.StartedAt,
		&workshift.FinishedAt,
		&workshift.Duration,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, workshifts.ErrWorkshiftNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &workshift, nil
}

type WorkshiftSaver struct {
	db *sql.DB
}

func NewWorkshiftSaver(db *sql.DB) *WorkshiftSaver {
	return &WorkshiftSaver{db: db}
}

// Save implements [commands.StartWorkshiftRepositoryPort].
func (ws *WorkshiftSaver) Save(ctx context.Context, workshift *workshifts.Workshift) error {
	const op = "adapters.postgres.workshifts_repository.WorkshiftSaver.Save"
	if len(workshift.Events) == 0 {
		return nil
	}
	tx, err := ws.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range workshift.Events {
		switch e := event.(type) {
		case workshifts.StartedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO workshifts (id, name, status, started_at, finished_at, duration) VALUES ($1, $2, $3, $4, $5, $6);`,
				e.ID, e.Name, e.Status, e.StartedAt, workshift.FinishedAt, workshift.Duration,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case workshifts.FinishedEvent:
			if _, err := tx.ExecContext(ctx,
				`UPDATE workshifts SET status = $1, finished_at = $2, duration = $3 WHERE id = $4;`,
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

type WorkshiftsRepository struct {
	*WorkshiftActiveChecker
	*WorkshiftReader
	*WorkshiftSaver
}

func NewWorkshiftsRepository(
	activeChecker *WorkshiftActiveChecker,
	reader *WorkshiftReader,
	saver *WorkshiftSaver,
) *WorkshiftsRepository {
	return &WorkshiftsRepository{
		WorkshiftActiveChecker: activeChecker,
		WorkshiftReader:        reader,
		WorkshiftSaver:         saver,
	}
}
