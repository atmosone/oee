package openapi

import (
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/equipment/commands"
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

func (h *Handler) CreateEquipment(w http.ResponseWriter, r *http.Request) {
	const op = "app.equipment.adapters.openapi.Handler.CreateEquipment"
	if err := h.app.HandleCommand(r.Context(), commands.CreateCommand{}); err != nil {
		switch {

		}
	}
}
