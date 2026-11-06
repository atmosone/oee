package openapi

import (
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/qualification/commands"
)

type App interface {
	cqrs.CommandHandler
}

type Handler struct {
	app    App
	logger *slog.Logger
}

func New(app App, logger *slog.Logger) *Handler {
	return &Handler{logger: logger, app: app}
}

func (h *Handler) CreateQualificationRule(w http.ResponseWriter, r *http.Request) {
	const op = "app.qualification.adapters.openapi.Handler.CreateQualificationRule"
	if err := h.app.HandleCommand(r.Context(), commands.CreateRuleCommand{}); err != nil {
		switch {

		}
	}
}

func (h *Handler) RemoveQualificationRule(w http.ResponseWriter, r *http.Request) {
	const op = "app.qualification.adapters.openapi.Handler.RemoveQualificationRule"
	if err := h.app.HandleCommand(r.Context(), commands.RemoveRuleCommand{}); err != nil {
		switch {

		}
	}
}

func (h *Handler) QualifyDowntime(w http.ResponseWriter, r *http.Request) {
	const op = "app.qualification.adapters.openapi.Handler.QualifyDowntime"
	if err := h.app.HandleCommand(r.Context(), commands.QualifyDowntimeCommand{}); err != nil {
		switch {

		}
	}
}
