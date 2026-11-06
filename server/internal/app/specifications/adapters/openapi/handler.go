package openapi

import (
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/specifications/commands"
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

func (h *Handler) CreateSpecification(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.CreateCommand{}); err != nil {
		switch {
		}
	}
}

func (h *Handler) DeprecateSpecification(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.DeprecateCommand{}); err != nil {
		switch {
		}
	}
}
