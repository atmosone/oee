package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/intervals"
	"github.com/atmosone/oee/internal/domain/production"
	"github.com/atmosone/oee/internal/domain/workshifts"
	"github.com/atmosone/oee/internal/views/analytics"
)

type WorkshiftsIntervalsProjection struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewWorkshiftsIntervalsProjection(
	logger *slog.Logger,
	db *sql.DB,
) *WorkshiftsIntervalsProjection {
	return &WorkshiftsIntervalsProjection{
		logger: logger,
		db:     db,
	}
}

func (wip *WorkshiftsIntervalsProjection) Types() []cqrs.EventType {
	return []cqrs.EventType{
		workshifts.EventTypeStarted,
		intervals.EventTypeStarted,
	}
}

func (wip *WorkshiftsIntervalsProjection) HandleEvent(ctx context.Context, event cqrs.Event) {
	const op = "adapters.postgres.projections.WorkshiftsIntervalsProjection.HandleEvent"
	switch e := event.(type) {
	case workshifts.StartedEvent:
		if _, err := wip.db.ExecContext(ctx,
			`INSERT INTO workshifts_intervals (workshift_id, interval_id, equipment_id)
				SELECT $1, id, equipment_id FROM intervals WHERE status = $2
				ON CONFLICT DO NOTHING;`,
			e.ID, intervals.StatusActive,
		); err != nil {
			wip.logger.ErrorContext(ctx, "failed to create relations",
				slog.String("op", op),
				slog.String("event_type", e.Type().String()),
				slog.String("workshift_id", e.ID),
				slog.Any("error", err),
			)
		}
	case intervals.StartedEvent:
		if _, err := wip.db.ExecContext(ctx,
			`INSERT INTO workshifts_intervals (workshift_id, interval_id, equipment_id)
					SELECT id, $1, $2 FROM workshifts WHERE status = $3
					ON CONFLICT DO NOTHING;`,
			e.ID, e.EquipmentID, workshifts.StatusActive,
		); err != nil {
			wip.logger.ErrorContext(ctx, "failed to create relations",
				slog.String("op", op),
				slog.String("event_type", e.Type().String()),
				slog.String("interval_id", e.ID),
				slog.String("equipment_id", e.EquipmentID),
				slog.Any("error", err),
			)
		}
	}
}

type WorkshiftProductionCounter struct {
	logger *slog.Logger
	db     *sql.DB
}

// HandleEvent implements [cqrs.Projection].
func (wpc *WorkshiftProductionCounter) HandleEvent(ctx context.Context, e cqrs.Event) {
	const op = "adapters.postgres.projections.WorkshiftProductionCounter.HandleEvent"
	event := e.(production.IncreasedEvent)
	var query string
	switch production.Category(event.Category) {
	case production.CategoryGood:
		query = `INSERT INTO workshifts_production (workshift_id, equipment_id, good)
		SELECT id, $1, $2 FROM workshifts WHERE status = $3
		ON CONFLICT (workshift_id, equipment_id) DO UPDATE SET
    	good = workshifts_production.good + EXCLUDED.good;`
	case production.CategoryDefect:
		query = `INSERT INTO workshifts_production (workshift_id, equipment_id, defect)
		SELECT id, $1, $2 FROM workshifts WHERE status = $3
		ON CONFLICT (workshift_id, equipment_id) DO UPDATE SET
    	defect = workshifts_production.defect + EXCLUDED.defect;`
	default:
		return
	}
	if _, err := wpc.db.ExecContext(ctx,
		query,
		event.EquipmentID, event.Growth, workshifts.StatusActive,
	); err != nil {
		wpc.logger.ErrorContext(ctx, "failed to update workshift production",
			slog.Any("error", err))
	}
}

// Types implements [cqrs.Projection].
func (wpc *WorkshiftProductionCounter) Types() []cqrs.EventType {
	return []cqrs.EventType{production.EventTypeIncreased}
}

func NewWorkshiftProductionCounter(
	logger *slog.Logger,
	db *sql.DB,
) *WorkshiftProductionCounter {
	return &WorkshiftProductionCounter{
		logger: logger,
		db:     db,
	}
}

type OeeMetricsBuilder struct {
	logger *slog.Logger
	db     *sql.DB
}

// HandleEvent implements [cqrs.Projection].
func (omb *OeeMetricsBuilder) HandleEvent(ctx context.Context, e cqrs.Event) {
	const op = ""
	go func(ctx context.Context, e workshifts.FinishedEvent) {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		workshift, err := omb.workshift(ctx, e.ID)
		if err != nil {
			return
		}
		equipment, err := omb.equipment(ctx)
		if err != nil {
			return
		}
		for _, eq := range equipment {
			production, err := omb.production(ctx, e.ID, eq.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return
			}
			downtimes, err := omb.downtimes(ctx, e.ID, eq.ID)
			if err != nil {
				return
			}
			availability := analytics.ComputeWorkshiftEquipmentAvailabilityRate(
				e.FinishedAt, workshift, downtimes,
			)
			performance := analytics.ComputeWorkshiftEquipmentPerformanceRate()
			quality := analytics.ComputeWorkshiftEquipmentQualityRate(production)
			oee := analytics.ComputeWorkshiftEquipmentOEE(availability, performance, quality)
			if _, err := omb.db.ExecContext(ctx,
				`INSERT INTO oee_metrics (workshift_id, equipment_id, availability, performance, quality, oee) 
				VALUES ($1, $2, $3, $4, $5, $6);`,
				workshift.ID, eq.ID, availability, performance, quality, oee,
			); err != nil {
				return
			}
		}
	}(context.WithoutCancel(ctx), e.(workshifts.FinishedEvent))
}

func (omb *OeeMetricsBuilder) workshift(ctx context.Context,
	workshiftID string) (analytics.Workshift, error) {
	var res analytics.Workshift
	if err := omb.db.QueryRowContext(ctx,
		`SELECT id, name, started_at, finished_at FROM workshifts WHERE id = $1;`,
		workshiftID,
	).Scan(
		&res.ID, &res.Name, &res.Interval.StartedAt, &res.Interval.FinishedAt,
	); err != nil {
		return res, err
	}
	return res, nil
}

func (omb *OeeMetricsBuilder) production(ctx context.Context,
	workshiftID, equipmentID string) (analytics.Production, error) {
	var res analytics.Production
	if err := omb.db.QueryRowContext(ctx,
		`SELECT good, defect 
		FROM workshifts_production 
		WHERE workshift_id = $1 
		AND equipment_id = $2;`,
		workshiftID, equipmentID,
	).Scan(&res.Good, &res.Defect); err != nil {
		return analytics.Production{}, err
	}
	return res, nil
}

func (omb *OeeMetricsBuilder) equipment(ctx context.Context) ([]analytics.Equipment, error) {
	list := make([]analytics.Equipment, 0)
	rows, err := omb.db.QueryContext(ctx, `SELECT id, name, category FROM equipment;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item analytics.Equipment
		if err := rows.Scan(&item.ID, &item.Category, &item.Name); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (omb *OeeMetricsBuilder) downtimes(ctx context.Context,
	workshiftID, equipmentID string) ([]analytics.DowntimeInterval, error) {
	list := make([]analytics.DowntimeInterval, 0)
	rows, err := omb.db.QueryContext(ctx,
		`SELECT i.started_at, i.finished_at FROM intervals i
		JOIN workshifts_intervals wi ON i.id = wi.interval_id
		WHERE wi.workshift_id = $1 AND wi.equipment_id = $2 AND i.category = 'downtime';`,
		workshiftID, equipmentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item analytics.DowntimeInterval
		if err := rows.Scan(
			&item.Interval.StartedAt,
			&item.Interval.FinishedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

// Types implements [cqrs.Projection].
func (o *OeeMetricsBuilder) Types() []cqrs.EventType {
	return []cqrs.EventType{workshifts.EventTypeFinished}
}

func NewOeeMetricsBuilder(
	logger *slog.Logger,
	db *sql.DB,
) *OeeMetricsBuilder {
	return &OeeMetricsBuilder{
		logger: logger,
		db:     db,
	}
}
