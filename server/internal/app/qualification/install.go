package qualification

import (
	"database/sql"
	"log/slog"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/qualification/adapters/openapi"
	"github.com/atmosone/oee/internal/app/qualification/adapters/postgres"
	"github.com/atmosone/oee/internal/app/qualification/commands"
)

func Install(app cqrs.App, db *sql.DB, logger *slog.Logger) *openapi.Handler {
	repository := postgres.New(db)
	app.OnCommand(commands.CommandTypeCreateRule, commands.CreateRule(repository))
	app.OnCommand(commands.CommandTypeRemoveRule, commands.RemoveRule(repository))
	app.OnCommand(commands.CommandTypeQualifyDowntime, commands.QualifyDowntime(repository))
	return openapi.New(app, logger)
}
