package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/intervals"
	"github.com/atmosone/oee/internal/domain/workshifts"
	"github.com/atmosone/oee/internal/queries"
	"github.com/atmosone/oee/internal/views"
	"github.com/atmosone/oee/pkg/openapi"
)

func LinesList(db *sql.DB) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		const op = "adapters.postgres.queries.LinesList"
		res := make(openapi.LinesList, 0)
		rows, err := db.QueryContext(ctx,
			`SELECT id, name FROM lines;`,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
		for rows.Next() {
			var item openapi.LinesListItem
			if err := rows.Scan(&item.Id, &item.Name); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			res = append(res, item)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return res, nil
	}
}

func LineTopology(db *sql.DB) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		const op = "adapters.postgres.queries.LineTopology"
		query := q.(queries.LineTopologyQuery)
		res := make(openapi.LineTopology, 0)
		rows, err := db.QueryContext(ctx,
			`SELECT e.id, e.name, e.category FROM equipment e 
			JOIN lines_equipment le ON e.id = le.equipment_id
			WHERE le.line_id = $1
			ORDER BY le.position ASC;`,
			query.Id,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
		for rows.Next() {
			var item openapi.LineTopologyItem
			if err := rows.Scan(&item.Id, &item.Name, &item.Category); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			res = append(res, item)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return res, nil
	}
}

func CurrentWorkshift(db *sql.DB) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		const op = "adapters.postgres.queries.CurrentWorkshift"
		var res openapi.CurrentWorkshift
		if err := db.QueryRowContext(ctx,
			`SELECT id, name, started_at FROM workshifts WHERE status = $1;`,
			workshifts.StatusActive,
		).Scan(&res.Id, &res.Name, &res.StartedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, queries.ErrNotFound
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return res, nil
	}
}

func CurrentInterval(db *sql.DB) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		const op = "adapters.postgres.queries.CurrentInterval"
		query := q.(queries.CurrentIntervalQuery)
		var res openapi.CurrentInterval
		if err := db.QueryRowContext(ctx,
			`SELECT id, equipment_id, category, started_at 
			FROM intervals 
			WHERE equipment_id = $1 AND status = $2;`,
			query.EquipmentId, intervals.StatusActive,
		).Scan(
			&res.Id, &res.EquipmentId, &res.Category, &res.StartedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, queries.ErrNotFound
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return res, nil
	}
}

func CurrentIntervalsList(db *sql.DB) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		const op = "adapters.postgres.queries.CurrentIntervalsList"
		res := make(openapi.CurrentIntervalsList, 0)
		rows, err := db.QueryContext(ctx,
			`SELECT id, equipment_id, category, started_at FROM intervals WHERE status = $1;`,
			intervals.StatusActive,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
		for rows.Next() {
			var item openapi.CurrentInterval
			if err := rows.Scan(
				&item.Id, &item.EquipmentId, &item.Category, &item.StartedAt,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			res = append(res, item)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return res, nil
	}
}

func EquipmentTimeline(db *sql.DB) cqrs.QueryHandlerFunc {
	return func(ctx context.Context, q cqrs.Query) (cqrs.Result, error) {
		const op = "adapters.postgres.queries.EquipmentTimeline"
		query := q.(queries.EquipmentTimelineQuery)
		var res views.Timeline
		if err := db.QueryRowContext(ctx,
			`SELECT id, name FROM equipment WHERE id = $1;`,
			query.EquipmentID,
		).Scan(&res.Equipment.ID,
			&res.Equipment.Name,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%w (equipment not found)", queries.ErrNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if err := db.QueryRowContext(ctx,
			`SELECT id, started_at, finished_at FROM workshifts WHERE id = $1;`,
			query.WorkshiftID,
		).Scan(&res.Workshift.ID,
			&res.Workshift.Interval.StartedAt,
			&res.Workshift.Interval.FinishedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%w (workshift not found)", queries.ErrNotFound)
			}
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rows, err := db.QueryContext(ctx,
			`SELECT i.id, i.category, i.started_at, i.finished_at FROM intervals i
			JOIN workshifts_intervals wi ON i.id = wi.interval_id
			WHERE wi.workshift_id = $1 
			AND i.equipment_id = $2
			ORDER BY i.started_at ASC;`,
			query.WorkshiftID, query.EquipmentID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		defer rows.Close()
		for rows.Next() {
			var item views.TimelineRecord
			if err := rows.Scan(&item.ID,
				&item.Category,
				&item.Interval.StartedAt,
				&item.Interval.FinishedAt,
			); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			res.Records = append(res.Records, item)
		}
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return res, nil
	}
}
