package equipment

import (
	"database/sql"
	"log/slog"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/equipment/adapters/openapi"
	"github.com/atmosone/oee/internal/app/equipment/adapters/postgres"
	"github.com/atmosone/oee/internal/app/equipment/commands"
)

func Install(app cqrs.App, db *sql.DB, logger *slog.Logger) *openapi.Handler {
	repository := postgres.New(db)
	app.OnCommand(commands.CommandTypeCreate, commands.Create(repository))
	return openapi.New(app, logger)
}
