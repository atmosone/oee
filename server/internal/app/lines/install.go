package lines

import (
	"database/sql"
	"log/slog"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/lines/adapters/openapi"
	"github.com/atmosone/oee/internal/app/lines/adapters/postgres"
	"github.com/atmosone/oee/internal/app/lines/commands"
)

func Install(app cqrs.App, db *sql.DB, logger *slog.Logger) *openapi.Handler {
	repository := postgres.New(db)
	app.OnCommand(commands.CommandTypeCreate, commands.Create(repository))
	app.OnCommand(commands.CommandTypeAddEquipment, commands.AddEquipment(repository))
	app.OnCommand(commands.CommandTypeMoveEquipment, commands.MoveEquipment(repository))
	app.OnCommand(commands.CommandTypeRemoveEquipment, commands.RemoveEquipment(repository))
	return openapi.New(app, logger)
}
