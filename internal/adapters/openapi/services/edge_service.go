package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/commands"
	"github.com/atmosone/oee/internal/domain/intervals"
	"github.com/atmosone/oee/internal/domain/production"
	"github.com/atmosone/oee/internal/domain/qualification"
)

type EdgeService struct {
	logger *slog.Logger
	app    cqrs.CommandHandler
}

func NewEdgeService(logger *slog.Logger, app cqrs.CommandHandler) *EdgeService {
	return &EdgeService{logger: logger, app: app}
}

// StartDowntime implements [openapi.ServerInterface].
func (s *EdgeService) FinishInterval(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EdgeService.FinishInterval"
	var cmd commands.FinishIntervalCommand
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		s.logger.DebugContext(ctx, "failed to decode command",
			slog.String("op", op),
			slog.String("command_type", cmd.Type().String()),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.app.HandleCommand(ctx, cmd); err != nil {
		switch {
		case errors.Is(err, intervals.ErrNotFound):
			s.logger.DebugContext(ctx, "failed to finish interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, intervals.ErrValidationFailed):
			s.logger.DebugContext(ctx, "failed to finish interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, intervals.ErrAlreadyFinished):
			s.logger.DebugContext(ctx, "failed to finish interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to finish interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (s *EdgeService) StartInterval(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EdgeService.StartInterval"
	var cmd commands.StartIntervalCommand
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		s.logger.DebugContext(ctx, "failed to decode command",
			slog.String("op", op),
			slog.String("command_type", cmd.Type().String()),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.app.HandleCommand(ctx, cmd); err != nil {
		switch {
		case errors.Is(err, intervals.ErrValidationFailed):
			s.logger.DebugContext(ctx, "failed to start interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, intervals.ErrAlreadyExists):
			s.logger.DebugContext(ctx, "failed to start interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to start interval",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// IncreaseProduction implements [openapi.ServerInterface].
func (s *EdgeService) IncreaseProduction(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EdgeService.IncreaseProduction"
	var cmd commands.IncreaseProductionCommand
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		s.logger.DebugContext(ctx, "failed to decode command",
			slog.String("op", op),
			slog.String("command_type", cmd.Type().String()),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.app.HandleCommand(ctx, cmd); err != nil {
		switch {
		case errors.Is(err, production.ErrActiveIntervalNotFound):
			s.logger.DebugContext(ctx, "failed to increase production",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, production.ErrValidationFailed):
			s.logger.DebugContext(ctx, "failed to increase production",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to increase production",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

func (s *EdgeService) QualifyDowntime(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EdgeService.QualifyDowntime"
	var cmd commands.QualifyDowntimeCommand
	ctx := r.Context()
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		s.logger.DebugContext(ctx, "failed to decode command",
			slog.String("op", op),
			slog.String("command_type", cmd.Type().String()),
			slog.Any("error", err),
		)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.app.HandleCommand(ctx, cmd); err != nil {
		switch {
		case errors.Is(err, qualification.ErrDowntimeNotFound):
			s.logger.DebugContext(ctx, "failed to qualify downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, qualification.ErrEquipmentConflict):
			s.logger.DebugContext(ctx, "failed to qualify downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		case errors.Is(err, qualification.ErrForbiddenSpecification) ||
			errors.Is(err, qualification.ErrDeprecatedSpecification) ||
			errors.Is(err, qualification.ErrDowntimeNotFinishedYet):
			s.logger.DebugContext(ctx, "failed to qualify downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to start downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
