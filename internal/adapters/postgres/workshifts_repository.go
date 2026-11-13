package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/atmosone/oee/internal/domain/workshifts"
)

type WorkshiftsReader struct {
	db *sql.DB
}

func NewWorkshiftsReader(db *sql.DB) *WorkshiftsReader {
	return &WorkshiftsReader{db: db}
}

func (wr *WorkshiftsReader) Changeover(ctx context.Context) (*workshifts.Changeover, error) {
	const op = "adapters.postgres.workshifts_repository.WorkshiftsReader.Changeover"
	var changeover workshifts.Changeover
	var active workshifts.Workshift
	if err := wr.db.QueryRowContext(ctx,
		`SELECT id, name, status, started_at FROM workshifts WHERE status = $1;`,
		workshifts.StatusActive,
	).Scan(&active.ID,
		&active.Name,
		&active.Status,
		&active.StartedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &changeover, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	changeover.Active = &active
	return &changeover, nil
}

type WorkshiftsSaver struct {
	db *sql.DB
}

func NewWorkshiftsSaver(db *sql.DB) *WorkshiftsSaver {
	return &WorkshiftsSaver{db: db}
}

// Save implements [commands.StartWorkshiftRepositoryPort].
func (ws *WorkshiftsSaver) Save(ctx context.Context, changeover *workshifts.Changeover) error {
	const op = "adapters.postgres.workshifts_repository.WorkshiftsSaver.Save"
	if len(changeover.Events) == 0 {
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
	for _, event := range changeover.Events {
		switch e := event.(type) {
		case workshifts.StartedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO workshifts (id, name, status, started_at, finished_at, duration) 
				VALUES ($1, $2, $3, $4, $5, $6);`,
				e.ID, e.Name, e.Status, e.StartedAt, time.Time{}, 0,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case workshifts.FinishedEvent:
			if _, err := tx.ExecContext(ctx,
				`UPDATE workshifts 
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

type WorkshiftsRepository struct {
	*WorkshiftsReader
	*WorkshiftsSaver
}

func NewWorkshiftsRepository(
	reader *WorkshiftsReader,
	saver *WorkshiftsSaver,
) *WorkshiftsRepository {
	return &WorkshiftsRepository{
		WorkshiftsReader: reader,
		WorkshiftsSaver:  saver,
	}
}
