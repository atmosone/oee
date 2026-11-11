package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/commands"
	"github.com/atmosone/oee/internal/domain/downtimes"
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
func (s *EdgeService) StartDowntime(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EdgeService.StartDowntime"
	var cmd commands.StartDowntimeCommand
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
		case errors.Is(err, downtimes.ErrInvalidStartTime):
			s.logger.DebugContext(ctx, "failed to finish downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, downtimes.ErrActiveAlreadyExists):
			s.logger.DebugContext(ctx, "failed to start downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
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

// FinishDowntime implements [openapi.ServerInterface].
func (s *EdgeService) FinishDowntime(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EdgeService.FinishDowntime"
	var cmd commands.FinishDowntimeCommand
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
		case errors.Is(err, downtimes.ErrDowntimeNotFound):
			s.logger.DebugContext(ctx, "failed to finish downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, downtimes.ErrInvalidFinishTime):
			s.logger.DebugContext(ctx, "failed to finish downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, downtimes.ErrAlreadyFinished):
			s.logger.DebugContext(ctx, "failed to finish downtime",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
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

// QualifyDowntime implements [openapi.ServerInterface].
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
