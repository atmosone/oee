package specifications

import (
	"database/sql"
	"log/slog"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/specifications/adapters/openapi"
	"github.com/atmosone/oee/internal/app/specifications/adapters/postgres"
	"github.com/atmosone/oee/internal/app/specifications/commands"
)

func Install(app cqrs.App, db *sql.DB, logger *slog.Logger) *openapi.Handler {
	repository := postgres.New(db)
	app.OnCommand(commands.CommandTypeCreate, commands.Create(repository))
	app.OnCommand(commands.CommandTypeDeprecate, commands.Deprecate(repository))
	return openapi.New(app, logger)
}
