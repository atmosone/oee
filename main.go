package main

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	_ "embed"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/adapters/openapi"
	"github.com/atmosone/oee/internal/adapters/openapi/services"
	"github.com/atmosone/oee/internal/adapters/postgres"
	"github.com/atmosone/oee/internal/commands"
	"github.com/atmosone/oee/internal/queries"
	_ "github.com/lib/pq"
)

const DSN = "postgres://postgres:postgres@localhost:5432/oee?sslmode=disable"

//go:embed sql/postgres/schema.sql
var schema string

func migrate(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, schema); err != nil {
		return err
	}
	return tx.Commit()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM)
	defer cancel()
	// Install logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	// Application Data Access Layer
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		logger.ErrorContext(ctx, "failed to ping database",
			slog.Any("error", err))
		return
	}
	if err := migrate(ctx, db); err != nil {
		logger.ErrorContext(ctx, "failed to migrate schema",
			slog.Any("error", err))
		return
	}
	// Install equipment repository
	equipment := postgres.NewEquipmentRepository(
		postgres.NewEquipmentSaver(db),
	)
	// Install lines repository
	lines := postgres.NewLinesRepository(
		postgres.NewLineReader(db),
		postgres.NewLineSaver(db),
	)
	// Install worktimes repository
	intervals := postgres.NewIntervalsRepository(
		postgres.NewIntervalReader(db),
		postgres.NewIntervalSaver(db),
	)
	// Install specifications repository
	specifications := postgres.NewSpecificationsRepository(
		postgres.NewSpecificationReader(db),
		postgres.NewSpecificationSaver(db),
	)
	// Install workshifts repository
	workshifts := postgres.NewWorkshiftsRepository(
		postgres.NewWorkshiftsReader(db),
		postgres.NewWorkshiftsSaver(db),
	)
	qualification := postgres.NewQualificationRepository(
		postgres.NewQualificationSpecificationReader(db),
		postgres.NewQualificationDowntimeReader(db),
		postgres.NewQualificationReader(db),
		postgres.NewQualificationSaver(db),
	)
	production := postgres.NewProductionRepository(
		postgres.NewProductionReader(db),
		postgres.NewProductionSaver(db),
	)
	// Install [cqrs.App]
	app := cqrs.New()
	app.OnPanic(func(r any) { debug.PrintStack() })
	// Command side
	app.OnCommand(commands.CommandTypeCreateEquipment, commands.CreateEquipment(equipment))
	app.OnCommand(commands.CommandTypeCreateLine, commands.CreateLine(lines))
	app.OnCommand(commands.CommandTypeAddEquipmentToLine, commands.AddEquipmentToLine(lines))
	app.OnCommand(commands.CommandTypeMoveEquipmentInLine, commands.MoveEquipmentInLine(lines))
	app.OnCommand(commands.CommandTypeRemoveEquipmentFromLine, commands.RemoveEquipmentFromLine(lines))
	app.OnCommand(commands.CommandTypeStartInterval, commands.StartInterval(intervals))
	app.OnCommand(commands.CommandTypeFinishInterval, commands.FinishInterval(intervals))
	app.OnCommand(commands.CommandTypeQualifyDowntime, commands.QualifyDowntime(qualification))
	app.OnCommand(commands.CommandTypeCreateQualificationRule, commands.CreateQualificationRule(qualification))
	app.OnCommand(commands.CommandTypeRemoveQualificationRule, commands.RemoveQualificationRule(qualification))
	app.OnCommand(commands.CommandTypeCreateSpecification, commands.CreateSpecification(specifications))
	app.OnCommand(commands.CommandTypeDeprecateSpecification, commands.DeprecateSpecification(specifications))
	app.OnCommand(commands.CommandTypeStartWorkshift, commands.StartWorkshift(workshifts))
	app.OnCommand(commands.CommandTypeFinishWorkshift, commands.FinishWorkshift(workshifts))
	app.OnCommand(commands.CommandTypeIncreaseProduction, commands.IncreaseProduction(production))
	// Query side
	app.OnQuery(queries.QueryTypeLinesList, postgres.LinesList(db))
	app.OnQuery(queries.QueryTypeLineTopology, postgres.LineTopology(db))
	app.OnQuery(queries.QueryTypeCurrentInterval, postgres.CurrentInterval(db))
	app.OnQuery(queries.QueryTypeCurrentWorkshift, postgres.CurrentWorkshift(db))
	app.OnQuery(queries.QueryTypeCurrentIntervalsList, postgres.CurrentIntervalsList(db))
	app.Build(postgres.NewWorkshiftsIntervalsProjection(logger, db))
	app.Build(postgres.NewWorkshiftProductionCounter(logger, db))
	app.Build(postgres.NewOeeMetricsBuilder(logger, db))

	router := http.NewServeMux()
	server := &http.Server{
		Addr: ":8880",
		Handler: openapi.New(router,
			services.NewManagementService(logger, app),
			services.NewEdgeService(logger, app),
			services.NewEventService(logger, app),
			services.NewViewService(logger, app),
		),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			logger.ErrorContext(ctx, "server error",
				slog.Any("error", err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(), 5*time.Second,
	)
	defer shutdownCancel()
	if err := app.Close(); err != nil {
		logger.ErrorContext(ctx, "failed to close application",
			slog.Any("error", err))
	}
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.ErrorContext(ctx, "failed to close server",
			slog.Any("error", err))
	}
}
