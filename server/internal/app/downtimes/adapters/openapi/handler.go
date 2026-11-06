package openapi

import (
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/downtimes/commands"
)

type App interface {
	cqrs.CommandHandler
}

type Handler struct {
	app    App
	logger *slog.Logger
}

func New(app App, logger *slog.Logger) *Handler {
	return &Handler{app: app, logger: logger}
}

func (h *Handler) StartDowntime(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.StartCommand{}); err != nil {
		switch {

		}
	}
}

func (h *Handler) FinishDowntime(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.FinishCommand{}); err != nil {
		switch {

		}
	}
}
