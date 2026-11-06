package openapi

import (
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/app/lines/commands"
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

func (h *Handler) CreateLine(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.CreateCommand{}); err != nil {
		switch {

		}
	}
}

func (h *Handler) AddLineEquipment(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.AddEquipmentCommand{}); err != nil {
		switch {

		}
	}
}

func (h *Handler) MoveLineEquipment(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.MoveEquipmentCommand{}); err != nil {
		switch {

		}
	}
}

func (h *Handler) RemoveLineEquipment(w http.ResponseWriter, r *http.Request) {
	if err := h.app.HandleCommand(r.Context(), commands.RemoveEquipmentCommand{}); err != nil {
		switch {

		}
	}
}
