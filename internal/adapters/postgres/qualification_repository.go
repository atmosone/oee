package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
)

type QualificationSpecificationReader struct {
	db *sql.DB
}

func NewQualificationSpecificationReader(db *sql.DB) *QualificationSpecificationReader {
	return &QualificationSpecificationReader{db: db}
}
func (qsr *QualificationSpecificationReader) Specification(ctx context.Context, id specifications.ID) (qualification.Specification, error) {
	const op = "adapters.postgres.qualification_repository.QualificationSpecificationReader.Specification"
	var specification qualification.Specification
	if err := qsr.db.QueryRowContext(ctx,
		`SELECT id, deprecated FROM specifications WHERE id = $1;`,
		id,
	).Scan(&specification.ID, &specification.Deprecated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return qualification.Specification{}, qualification.ErrSpecificationNotFound
		}
		return qualification.Specification{}, fmt.Errorf("%s: %w", op, err)
	}
	return specification, nil
}

type QualificationDowntimeReader struct {
	db *sql.DB
}

func NewQualificationDowntimeReader(db *sql.DB) *QualificationDowntimeReader {
	return &QualificationDowntimeReader{db: db}
}
func (qdr *QualificationDowntimeReader) Downtime(ctx context.Context, id downtimes.ID) (qualification.Downtime, error) {
	const op = "adapters.postgres.qualification_repository.QualificationDowntimeReader.Downtime"
	var downtime qualification.Downtime
	if err := qdr.db.QueryRowContext(ctx,
		`SELECT id, equipment_id, status FROM downtimes WHERE id = $1;`,
		id,
	).Scan(&downtime.ID, &downtime.EquipmentID, &downtime.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return qualification.Downtime{}, qualification.ErrDowntimeNotFound
		}
		return qualification.Downtime{}, fmt.Errorf("%s: %w", op, err)
	}
	return downtime, nil
}

type QualificationReader struct {
	db *sql.DB
}

func NewQualificationReader(db *sql.DB) *QualificationReader {
	return &QualificationReader{db: db}
}
func (qr *QualificationReader) Qualification(ctx context.Context, id equipment.ID) (*qualification.Qualification, error) {
	const op = "adapters.postgres.qualification_repository.QualificationReader.Qualification"
	q := qualification.Qualification{
		EquipmentID:    id,
		Specifications: make(map[specifications.ID]qualification.Specification),
	}
	rows, err := qr.db.QueryContext(ctx,
		`SELECT FROM s.id, s.deprecated FROM specifications s
		JOIN qualification_rules qr ON qr.specification_id = s.id
		WHERE qr.equipment_id = $1;`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var specification qualification.Specification
		if err := rows.Scan(&specification.ID, &specification.Deprecated); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		q.Specifications[specification.ID] = specification
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &q, nil
}

type QualificationSaver struct {
	db *sql.DB
}

func NewQualificationSaver(db *sql.DB) *QualificationSaver {
	return &QualificationSaver{db: db}
}

// Save implements [commands.QualifyDowntimeRepositoryPort].
func (qs *QualificationSaver) Save(ctx context.Context, q *qualification.Qualification) error {
	const op = "adapters.postgres.qualification_repository.QualificationSaver.Save"
	if len(q.Events) == 0 {
		return nil
	}
	tx, err := qs.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range q.Events {
		switch e := event.(type) {
		case qualification.RuleCreatedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO qualification_rules (equipment_id, specification_id) VALUES ($1, $2);`,
				e.EquipmentID, e.SpecificationID,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case qualification.RuleRemovedEvent:
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM qualification_rules WHERE equipment_id = $1 AND specification_id = $2;`,
				e.EquipmentID, e.SpecificationID,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case qualification.DowntimeQualifiedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO downtimes_qualifications (downtime_id, specification_id, comment) VALUES ($1, $2, $3);`,
				e.DowntimeID, e.SpecificationID, e.Comment,
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

type QualificationRepository struct {
	*QualificationSpecificationReader
	*QualificationDowntimeReader
	*QualificationReader
	*QualificationSaver
}

func NewQualificationRepository(
	specificationReader *QualificationSpecificationReader,
	downtimeReader *QualificationDowntimeReader,
	reader *QualificationReader,
	saver *QualificationSaver,
) *QualificationRepository {
	return &QualificationRepository{
		QualificationSpecificationReader: specificationReader,
		QualificationDowntimeReader:      downtimeReader,
		QualificationReader:              reader,
		QualificationSaver:               saver,
	}
}
