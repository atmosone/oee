package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain/specifications"
)

type SpecificationReader struct{ db *sql.DB }

func NewSpecificationReader(db *sql.DB) *SpecificationReader {
	return &SpecificationReader{db: db}
}

// Specification implements [commands.DeprecateSpecificationRepositoryPort].
func (sr *SpecificationReader) Specification(ctx context.Context, id specifications.ID) (*specifications.Specification, error) {
	const op = "adapters.postgres.specifications_repository.SpecificationReader.Specification"
	var specification specifications.Specification
	if err := sr.db.QueryRowContext(ctx,
		`SELECT id, code, description, deprecated FROM specifications WHERE id = $1;`,
		id).Scan(&specification.ID,
		&specification.Code,
		&specification.Description,
		&specification.Deprecated,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, specifications.ErrSpecificationNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &specification, nil
}

type SpecificationSaver struct{ db *sql.DB }

func NewSpecificationSaver(db *sql.DB) *SpecificationSaver {
	return &SpecificationSaver{db: db}
}

// Save implements [commands.CreateSpecificationRepositoryPort].
func (ss *SpecificationSaver) Save(ctx context.Context, specification *specifications.Specification) error {
	const op = "adapters.postgres.specifications_repository.SpecificationSaver.Save"
	if len(specification.Events) == 0 {
		return nil
	}
	tx, err := ss.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range specification.Events {
		switch e := event.(type) {
		case specifications.CreatedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO specifications (id, code, description, deprecated) VALUES ($1, $2, $3, $4);`,
				e.ID, e.Code, e.Description, e.Deprecated,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case specifications.DeprecatedEvent:
			if _, err := tx.ExecContext(ctx,
				`UPDATE specifications SET deprecated = $1 WHERE id = $2;`,
				e.Deprecated, e.ID,
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

type SpecificationsRepository struct {
	*SpecificationReader
	*SpecificationSaver
}

func NewSpecificationsRepository(
	reader *SpecificationReader,
	saver *SpecificationSaver,
) *SpecificationsRepository {
	return &SpecificationsRepository{
		SpecificationReader: reader,
		SpecificationSaver:  saver,
	}
}
