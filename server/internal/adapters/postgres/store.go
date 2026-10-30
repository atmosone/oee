package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/atmosone/oee/internal/adapters/openapi/qualification"
)

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

// Equipment implements [queries.EquipmentStore].
func (s *Store) Equipment(
	ctx context.Context, params qualification.EquipmentParams,
) (*qualification.Equipment, error) {
	const op = "adapters.postgres.Store.Equipment"
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	var res qualification.Equipment
	if err := tx.QueryRowContext(ctx,
		`SELECT name FROM equipment WHERE id=$1;`,
		params.EquipmentId,
	).Scan(&res.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var line string
	if err := tx.QueryRowContext(ctx,
		`SELECT el.name FROM equipment_lines el
		JOIN equipment_line_items eli ON el.id = eli.equipment_line_id
		WHERE eli.equipment_id = $1;`,
		params.EquipmentId,
	).Scan(&line); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	} else {
		res.Line = &line
	}

	return &res, nil
}

// CurrentEquipmentWorkshift implements [queries.CurrentEquipmentWorkshiftStore].
func (s *Store) CurrentEquipmentWorkshift(ctx context.Context) (*qualification.CurrentEquipmentWorkshift, error) {
	const op = "adapters.postgres.Store.CurrentEquipmentWorkshift"
	var res qualification.CurrentEquipmentWorkshift
	if err := s.db.QueryRowContext(ctx,
		`SELECT id, name, started_at FROM workshifts 
		WHERE status = 'active' 
		ORDER BY started_at DESC 
		LIMIT 1;`,
	).Scan(&res.Id, &res.Name, &res.StartedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &res, nil
}

// CurrentEquipmentDowntime implements [queries.CurrentEquipmentDowntimeStore].
func (s *Store) CurrentEquipmentDowntime(
	ctx context.Context, params qualification.CurrentEquipmentDowntimeParams,
) (*qualification.CurrentEquipmentDowntime, error) {
	const op = "adapters.postgres.Store.CurrentEquipmentDowntime"
	var res qualification.CurrentEquipmentDowntime
	if err := s.db.QueryRowContext(ctx,
		`SELECT id, started_at FROM downtimes 
		WHERE equipment_id = $1 
		AND status = 'active';`,
		params.EquipmentId,
	).Scan(&res.Id, &res.StartedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &res, nil
}

// EquipmentSpecificationsList implements [queries.EquipmentSpecificationsListStore].
func (s *Store) EquipmentSpecificationsList(
	ctx context.Context, params qualification.EquipmentSpecificationsListParams,
) (qualification.EquipmentSpecificationsList, error) {
	const op = "adapters.postgres.Store.EquipmentSpecificationsList"
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx,
		`SELECT s.id, s.code, s.description FROM specifications s
		JOIN qualification_rules qr ON s.id = qr.specification_id
		WHERE qr.equipment_id = $1;`,
		params.EquipmentId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	list := make(qualification.EquipmentSpecificationsList, 0)
	for rows.Next() {
		var item qualification.EquipmentSpecificationsListItem
		if err := rows.Scan(&item.Id, &item.Code, &item.Description); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return list, nil
}

// UnqualifiedDowntimesList implements [queries.UnqualifiedDowntimesListStore].
func (s *Store) UnqualifiedDowntimesList(
	ctx context.Context, params qualification.UnqualifiedDowntimesListParams,
) (qualification.UnqualifiedDowntimesList, error) {
	const op = "adapters.postgres.Store.UnqualifiedDowntimesList"
	const threshold = 1 * time.Minute
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()
	rows, err := tx.QueryContext(ctx,
		`SELECT d.id, d.duration, d.started_at, d.finished_at FROM downtimes d
		WHERE NOT EXISTS (SELECT 1 FROM qualifications q WHERE q.downtime_id = d.id) 
		AND d.duration >= $1 
		AND d.equipment_id = $2 
		AND d.started_at > (SELECT w.started_at FROM workshifts w WHERE w.id = $3)
		ORDER BY d.started_at ASC;`,
		threshold, params.EquipmentId, params.WorkshiftId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()
	list := make(qualification.UnqualifiedDowntimesList, 0)
	for rows.Next() {
		var duration int64
		var item qualification.UnqualifiedDowntimesListItem
		if err := rows.Scan(&item.Id, &duration, &item.StartedAt, &item.FinishedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		item.Duration = time.Duration(duration).String()
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return list, nil
}
