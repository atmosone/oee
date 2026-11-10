package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/domain/downtimes"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

type WorkshiftDowntimesProjection struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewWorkshiftDowntimesProjection(
	logger *slog.Logger,
	db *sql.DB,
) *WorkshiftDowntimesProjection {
	return &WorkshiftDowntimesProjection{
		logger: logger,
		db:     db,
	}
}

func (wdp *WorkshiftDowntimesProjection) Types() []cqrs.EventType {
	return []cqrs.EventType{
		workshifts.EventTypeStarted,
		downtimes.EventTypeStarted,
	}
}

func (wdp *WorkshiftDowntimesProjection) HandleEvent(ctx context.Context, event cqrs.Event) {
	const op = "adapters.postgres.projections.WorkshiftDowntimesProjection.HandleEvent"
	switch e := event.(type) {
	case workshifts.StartedEvent:
		if _, err := wdp.db.ExecContext(ctx,
			`INSERT INTO workshifts_downtimes (workshift_id, downtime_id, equipment_id)
				SELECT $1, id, equipment_id FROM downtimes WHERE status = $2
				ON CONFLICT DO NOTHING;`,
			e.ID, downtimes.StatusActive,
		); err != nil {
			wdp.logger.ErrorContext(ctx, "failed to create relations",
				slog.String("op", op),
				slog.String("event_type", e.Type().String()),
				slog.String("workshift_id", e.ID),
				slog.Any("error", err),
			)
		}
	case downtimes.StartedEvent:
		if _, err := wdp.db.ExecContext(ctx,
			`INSERT INTO workshifts_downtimes (workshift_id, downtime_id, equipment_id)
					SELECT id, $1, $2 FROM workshifts WHERE status = $3
					ON CONFLICT DO NOTHING;`,
			e.ID, e.EquipmentID, workshifts.StatusActive,
		); err != nil {
			wdp.logger.ErrorContext(ctx, "failed to create relations",
				slog.String("op", op),
				slog.String("event_type", e.Type().String()),
				slog.String("downtime_id", e.ID),
				slog.String("equipment_id", e.EquipmentID),
				slog.Any("error", err),
			)
		}
	}

}
