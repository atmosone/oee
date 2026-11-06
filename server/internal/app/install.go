package app

import (
	"database/sql"
	"log/slog"
	"runtime/debug"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/downtimes"
	"github.com/atmosone/oee/internal/app/equipment"
	"github.com/atmosone/oee/internal/app/lines"
	"github.com/atmosone/oee/internal/app/qualification"
	"github.com/atmosone/oee/internal/app/specifications"
	"github.com/atmosone/oee/internal/app/workshifts"
)

func New(logger *slog.Logger, db *sql.DB) cqrs.App {
	app := cqrs.New()
	app.OnPanic(func(r any) {
		debug.PrintStack()
		logger.With(slog.Any("panic", r)).Error("recovered")
	})
	workshifts.Install(app, db, logger)
	equipment.Install(app, db, logger)
	lines.Install(app, db, logger)
	downtimes.Install(app, db, logger)
	specifications.Install(app, db, logger)
	qualification.Install(app, db, logger)
	return app
}
