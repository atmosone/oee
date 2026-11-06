package workshifts

import (
	"database/sql"
	"log/slog"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/workshifts/adapters/openapi"
	"github.com/atmosone/oee/internal/app/workshifts/adapters/postgres"
	"github.com/atmosone/oee/internal/app/workshifts/commands"
)

func Install(app cqrs.App, db *sql.DB, logger *slog.Logger) *openapi.Handler {
	repository := postgres.New(db)
	app.OnCommand(commands.CommandTypeStart, commands.Start(repository))
	app.OnCommand(commands.CommandTypeFinish, commands.Finish(repository))
	return openapi.New(app, logger)
}
