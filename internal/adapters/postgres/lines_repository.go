package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
)

type LineReader struct{ db *sql.DB }

func NewLineReader(db *sql.DB) *LineReader { return &LineReader{db: db} }

// Line implements [commands.AddEquipmentToLineRepositoryPort].
func (lr *LineReader) Line(ctx context.Context, id lines.ID) (*lines.Line, error) {
	const op = "adapters.postgres.lines_repository.LineReader.Line"
	tx, err := lr.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	var line lines.Line
	if err := tx.QueryRowContext(ctx,
		`SELECT id, name FROM lines WHERE id = $1;`,
		id,
	).Scan(&line.ID, &line.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, lines.ErrLineNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	rows, err := tx.QueryContext(ctx,
		`SELECT equipment_id FROM lines_equipment WHERE line_id = $1 ORDER BY position ASC;`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var equipmentID equipment.ID
		if err := rows.Scan(&equipmentID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		line.Equipment = append(line.Equipment, equipmentID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &line, nil
}

type LineSaver struct{ db *sql.DB }

func NewLineSaver(db *sql.DB) *LineSaver { return &LineSaver{db: db} }

// Save implements [commands.AddEquipmentToLineRepositoryPort].
func (ls *LineSaver) Save(ctx context.Context, line *lines.Line) error {
	const op = "adapters.postgres.lines_repository.LineSaver.Save"
	if len(line.Events) == 0 {
		return nil
	}
	tx, err := ls.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	for _, event := range line.Events {
		switch e := event.(type) {
		case lines.CreatedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO lines (id, name) VALUES ($1, $2);`,
				e.ID, e.Name,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case lines.EquipmentAddedEvent:
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO lines_equipment (line_id, equipment_id, position) VALUES ($1, $2, $3);`,
				e.ID, e.EquipmentID, e.Position,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case lines.EquipmentMovedEvent:
			if _, err := tx.ExecContext(ctx,
				`UPDATE lines_equipment SET position = $1 WHERE line_id = $2 AND equipment_id = $3;`,
				e.Position, e.ID, e.EquipmentID,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		case lines.EquipmentRemovedEvent:
			if _, err := tx.ExecContext(ctx,
				`DELETE FROM lines_equipment WHERE line_id = $1 AND equipment_id = $2;`,
				e.ID, e.EquipmentID,
			); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
		default:
			continue
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type LinesRepository struct {
	*LineReader
	*LineSaver
}

func NewLinesRepository(reader *LineReader, saver *LineSaver) *LinesRepository {
	return &LinesRepository{
		LineReader: reader,
		LineSaver:  saver,
	}
}
