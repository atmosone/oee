package postgres

import (
	"context"
	"database/sql"
	"errors"

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
		return nil, specifications.ErrSpecificationNotFound,
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
func (ss *SpecificationSaver) Save(context.Context, *specifications.Specification) error {
	panic("unimplemented")
}

type SpecificationsRepository struct {
	*SpecificationReader
	*SpecificationSaver
}

func NewSpecificationsRepository(
	reader *SpecificationReader
	saver *SpecificationSaver
) *SpecificationsRepository {
	return &SpecificationsRepository{
		SpecificationReader: reader,
SpecificationSaver: saver,
	}
}
