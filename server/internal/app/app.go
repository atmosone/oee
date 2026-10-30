package app

import (
	"database/sql"
	"log/slog"
	"runtime/debug"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/adapters/postgres"
	"github.com/atmosone/oee/internal/commands"
	"github.com/atmosone/oee/internal/queries"
)

func New(logger *slog.Logger, db *sql.DB) cqrs.App {
	app := cqrs.New()
	app.OnPanic(func(r any) {
		debug.PrintStack()
		logger.With(slog.Any("panic", r)).Error("recovered")
	})
	// Command Side
	repository := postgres.NewRepository(db)
	app.OnCommand(commands.CommandTypeStartWorkshift, commands.StartWorkshift(repository))
	app.OnCommand(commands.CommandTypeFinishWorkshift, commands.FinishWorkshift(repository))
	app.OnCommand(commands.CommandTypeCreateEquipment, commands.CreateEquipment(repository))
	app.OnCommand(commands.CommandTypeCreateEquipmentLine, commands.CreateEquipmentLine(repository))
	app.OnCommand(commands.CommandTypeExpandEquipmentLine, commands.ExpandEquipmentLine(repository))
	app.OnCommand(commands.CommandTypeCreateSpecification, commands.CreateSpecification(repository))
	app.OnCommand(commands.CommandTypeCreateQualificationRule, commands.CreateQualificationRule(repository))
	app.OnCommand(commands.CommandTypeStartDowntime, commands.StartDowntime(repository))
	app.OnCommand(commands.CommandTypeFinishDowntime, commands.FinishDowntime(repository))
	app.OnCommand(commands.CommandTypeCreateQualification, commands.CreateQualification(repository))
	// Query Side
	store := postgres.NewStore(db)
	app.OnQuery(queries.QueryTypeEquipment, queries.Equipment(store))
	app.OnQuery(queries.QueryTypeCurrentEquipmentDowntime, queries.CurrentEquipmentDowntime(store))
	app.OnQuery(queries.QueryTypeUnqualifiedDowntimesList, queries.UnqualifiedDowntimesList(store))
	app.OnQuery(queries.QueryTypeCurrentEquipmentWorkshift, queries.CurrentEquipmentWorkshift(store))
	app.OnQuery(queries.QueryTypeEquipmentSpecificationsList, queries.EquipmentSpecificationsList(store))
	return app
}
