package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/commands"
	"github.com/atmosone/oee/internal/domain/equipment"
	"github.com/atmosone/oee/internal/domain/lines"
	"github.com/atmosone/oee/internal/domain/qualification"
	"github.com/atmosone/oee/internal/domain/specifications"
	"github.com/atmosone/oee/internal/domain/workshifts"
)

type ManagementService struct {
	logger *slog.Logger
	app    cqrs.CommandHandler
}

func NewManagementService(logger *slog.Logger, app cqrs.CommandHandler) *ManagementService {
	return &ManagementService{logger: logger, app: app}
}

func (s *ManagementService) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.CreateEquipment"
	var cmd commands.CreateEquipmentCommand
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
		case errors.Is(err, equipment.ErrEquipmentAlreadyExists):
			s.logger.DebugContext(ctx, "failed to create equipment",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to create equipment",
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

func (s *ManagementService) CreateLine(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.CreateLine"
	var cmd commands.CreateLineCommand
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
		case errors.Is(err, lines.ErrLineAlreadyExists):
			s.logger.DebugContext(ctx, "failed to create line",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to create line",
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

// AddEquipmentToLine implements [openapi.ServerInterface].
func (s *ManagementService) AddEquipmentToLine(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.AddEquipmentToLine"
	var cmd commands.AddEquipmentToLineCommand
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
		case errors.Is(err, lines.ErrEquipmentAlreadyExists):
			s.logger.DebugContext(ctx, "failed to add equipment to line",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to add equipment to line",
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

// MoveEquipmentInLine implements [openapi.ServerInterface].
func (s *ManagementService) MoveEquipmentInLine(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.MoveEquipmentInLine"
	var cmd commands.MoveEquipmentInLineCommand
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
		case errors.Is(err, lines.ErrUnknownDirection):
			s.logger.DebugContext(ctx, "failed to move equipment in line",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, lines.ErrEquipmentNotExists):
			s.logger.DebugContext(ctx, "failed to move equipment in line",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, lines.ErrOutOfLine):
			s.logger.DebugContext(ctx, "failed to move equipment in line",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to move equipment in line",
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

// RemoveEquipmentFromLine implements [openapi.ServerInterface].
func (s *ManagementService) RemoveEquipmentFromLine(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.RemoveEquipmentFromLine"
	var cmd commands.RemoveEquipmentFromLineCommand
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
		case errors.Is(err, lines.ErrEquipmentNotExists):
			s.logger.DebugContext(ctx, "failed to remove equipment from line",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to remove equipment from line",
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

// CreateSpecification implements [openapi.ServerInterface].
func (s *ManagementService) CreateSpecification(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.CreateSpecification"
	var cmd commands.CreateSpecificationCommand
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
		case errors.Is(err, specifications.ErrAlreadyExists):
			s.logger.DebugContext(ctx, "failed to create specification",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to create specification",
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

// DeprecateSpecification implements [openapi.ServerInterface].
func (s *ManagementService) DeprecateSpecification(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.DeprecateSpecification"
	var cmd commands.DeprecateSpecificationCommand
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
		case errors.Is(err, specifications.ErrSpecificationNotFound):
			s.logger.DebugContext(ctx, "failed to deprecate specification",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, specifications.ErrAlreadyDeprecated):
			s.logger.DebugContext(ctx, "failed to deprecate specification",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to create specification",
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

// StartWorkshift implements [openapi.ServerInterface].
func (s *ManagementService) StartWorkshift(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.StartWorkshift"
	var cmd commands.StartWorkshiftCommand
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
		case errors.Is(err, workshifts.ErrValidationFailed):
			s.logger.DebugContext(ctx, "failed to start workshift",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to start workshift",
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

// FinishWorkshift implements [openapi.ServerInterface].
func (s *ManagementService) FinishWorkshift(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.FinishWorkshift"
	var cmd commands.FinishWorkshiftCommand
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
		case errors.Is(err, workshifts.ErrValidationFailed):
			s.logger.DebugContext(ctx, "failed to finish workshift",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, workshifts.ErrActiveNotExists):
			s.logger.DebugContext(ctx, "failed to start workshift",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to start workshift",
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

// CreateQualificationRule implements [openapi.ServerInterface].
func (s *ManagementService) CreateQualificationRule(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.CreateQualificationRule"
	var cmd commands.CreateQualificationRuleCommand
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
		case errors.Is(err, qualification.ErrSpecificationNotFound):
			s.logger.DebugContext(ctx, "failed to create qualification rule",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		case errors.Is(err, qualification.ErrDeprecatedSpecification):
			s.logger.DebugContext(ctx, "failed to create qualification rule",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		case errors.Is(err, qualification.ErrSpecificationAlreadyExists):
			s.logger.DebugContext(ctx, "failed to create qualification rule",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to start workshift",
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

// RemoveQualificationRule implements [openapi.ServerInterface].
func (s *ManagementService) RemoveQualificationRule(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.ManagementService.RemoveQualificationRule"
	var cmd commands.RemoveQualificationRuleCommand
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
		case errors.Is(err, qualification.ErrSpecificationNotExists):
			s.logger.DebugContext(ctx, "failed to remove qualification rule",
				slog.String("op", op),
				slog.String("command_type", cmd.Type().String()),
				slog.Any("error", err),
			)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			s.logger.ErrorContext(ctx, "failed to start workshift",
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
