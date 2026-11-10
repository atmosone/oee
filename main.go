package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/adapters/postgres"
	"github.com/atmosone/oee/internal/commands"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()
	// Install logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	db, err := sql.Open("postgres", "")
	if err != nil {
		return
	}
	defer db.Close()
	// Instell equipment repository
	equipment := postgres.NewEquipmentRepository(
		postgres.NewEquipmentSaver(db),
	)
	// Install lines repository
	lines := postgres.NewLinesRepository(
		postgres.NewLineReader(db),
		postgres.NewLineSaver(db),
	)
	// Install downtimes repository
	downtimes := postgres.NewDowntimesRepository(
		postgres.NewDowntimeActiveChecker(db),
		postgres.NewDowntimeReader(db),
		postgres.NewDowntimeSaver(db),
	)
	// Install specifications repository
	specifications := postgres.NewSpecificationsRepository(
		postgres.NewSpecificationReader(db),
		postgres.NewSpecificationSaver(db),
	)
	// Install [cqrs.App]
	app := cqrs.New()
	app.OnPanic(func(r any) { debug.PrintStack() })
	app.OnCommand(commands.CommandTypeCreateEquipment, commands.CreateEquipment(equipment))
	app.OnCommand(commands.CommandTypeCreateLine, commands.CreateLine(lines))
	app.OnCommand(commands.CommandTypeAddEquipmentToLine, commands.AddEquipmentToLine(lines))
	app.OnCommand(commands.CommandTypeMoveEquipmentInLine, commands.MoveEquipmentInLine(lines))
	app.OnCommand(commands.CommandTypeRemoveEquipmentFromLine, commands.RemoveEquipmentFromLine(lines))
	app.OnCommand(commands.CommandTypeStartDowntime, commands.StartDowntime(downtimes))
	app.OnCommand(commands.CommandTypeFinishDowntime, commands.FinishDowntime(downtimes))
	app.OnCommand(commands.CommandTypeQualifyDowntime, commands.QualifyDowntime())
	app.OnCommand(commands.CommandTypeCreateQualificationRule, commands.CreateQualificationRule())
	app.OnCommand(commands.CommandTypeRemoveQualificationRule, commands.RemoveQualificationRule())
	app.OnCommand(commands.CommandTypeCreateSpecification, commands.CreateSpecification(specifications))
	app.OnCommand(commands.CommandTypeDeprecateSpecification, commands.DeprecateSpecification(specifications))
	app.OnCommand(commands.CommandTypeStartWorkshift, commands.StartWorkshift())
	app.OnCommand(commands.CommandTypeFinishWorkshift, commands.FinishWorkshift())
	defer app.Close()
	<-ctx.Done()
}
