package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
)

type EventService struct {
	logger *slog.Logger
	app    cqrs.EventStreamer
}

func NewEventService(
	logger *slog.Logger,
	app cqrs.EventStreamer,
) *EventService {
	return &EventService{
		logger: logger,
		app:    app,
	}
}

func (s *EventService) Events(w http.ResponseWriter, r *http.Request) {
	const op = "adapters.openapi.services.EventService.Events"
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	stream, err := s.app.EventStream(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for event := range stream {
		data, err := json.Marshal(event)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\n", event.Type().String())
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}
