package openapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/adapters/openapi/edge"
	"github.com/atmosone/oee/internal/adapters/openapi/qualification"
	"github.com/atmosone/oee/internal/app/commands"
	"github.com/atmosone/oee/internal/app/queries"
	"github.com/atmosone/oee/internal/domain"
	"github.com/atmosone/oee/web"
)

type QualificationServerApp interface {
	cqrs.CommandHandler
	cqrs.QueryHandler
	cqrs.EventStreamer
}

type QualificationServer struct {
	*http.Server
	allowedEvents map[cqrs.EventType]struct{}
	app           QualificationServerApp
}

// Events implements [qualification.ServerInterface].
func (q *QualificationServer) Events(
	w http.ResponseWriter, r *http.Request, params qualification.EventsParams,
) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	stream, err := q.app.EventStream(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for event := range stream {
		fmt.Println(event.Type())
		if _, ok := q.allowedEvents[event.Type()]; !ok {
			continue
		}
		switch e := event.(type) {
		case domain.DowntimeStartedEvent:
			if e.EquipmentID != params.EquipmentId {
				continue
			}
		case domain.DowntimeFinishedEvent:
			if e.EquipmentID != params.EquipmentId {
				continue
			}
		case domain.EquipmentLineExpandedEvent:
			if e.EquipmentID != params.EquipmentId {
				continue
			}
		}
		data, err := json.Marshal(event)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\n", event.Type().String())
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}

// CreateQualification implements [qualification.ServerInterface].
func (q *QualificationServer) CreateQualification(w http.ResponseWriter, r *http.Request) {
	var cmd commands.CreateQualificationCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := q.app.HandleCommand(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// CurrentEquipmentDowntime implements [qualification.ServerInterface].
func (q *QualificationServer) CurrentEquipmentDowntime(w http.ResponseWriter, r *http.Request, params qualification.CurrentEquipmentDowntimeParams) {
	res, err := q.app.HandleQuery(r.Context(),
		queries.CurrentEquipmentDowntimeQuery{
			CurrentEquipmentDowntimeParams: params,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// CurrentEquipmentWorkshift implements [qualification.ServerInterface].
func (q *QualificationServer) CurrentEquipmentWorkshift(w http.ResponseWriter, r *http.Request) {
	res, err := q.app.HandleQuery(r.Context(),
		queries.CurrentEquipmentWorkshiftQuery{},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Equipment implements [qualification.ServerInterface].
func (q *QualificationServer) Equipment(w http.ResponseWriter, r *http.Request, params qualification.EquipmentParams) {
	res, err := q.app.HandleQuery(r.Context(),
		queries.EquipmentQuery{
			EquipmentParams: params,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// UnqualifiedDowntimesList implements [qualification.ServerInterface].
func (q *QualificationServer) UnqualifiedDowntimesList(
	w http.ResponseWriter, r *http.Request, params qualification.UnqualifiedDowntimesListParams,
) {
	res, err := q.app.HandleQuery(r.Context(),
		queries.UnqualifiedDowntimesListQuery{
			UnqualifiedDowntimesListParams: params,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// EquipmentSpecificationsList implements [qualification.ServerInterface].
func (q *QualificationServer) EquipmentSpecificationsList(
	w http.ResponseWriter, r *http.Request, params qualification.EquipmentSpecificationsListParams,
) {
	res, err := q.app.HandleQuery(r.Context(),
		queries.EquipmentSpecificationsListQuery{
			EquipmentSpecificationsListParams: params,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func NewQualificationServer(addr string, app QualificationServerApp) *QualificationServer {
	srv := &QualificationServer{
		Server: &http.Server{Addr: addr},
		app:    app,
		allowedEvents: map[cqrs.EventType]struct{}{
			domain.EventTypeWorkshiftStarted:      {},
			domain.EventTypeWorkshiftFinished:     {},
			domain.EventTypeDowntimeStarted:       {},
			domain.EventTypeDowntimeFinished:      {},
			domain.EventTypeEquipmentLineExpanded: {},
		},
	}
	mux := http.NewServeMux()
	mux.Handle("/", web.NewQualificationHandler())
	srv.Handler = qualification.HandlerWithOptions(srv, qualification.StdHTTPServerOptions{
		BaseRouter: mux,
		BaseURL:    "/api",
	})

	return srv
}

type EdgeApp interface {
	cqrs.EventStreamer
	cqrs.CommandHandler
}

type EdgeServer struct {
	*http.Server
	app           EdgeApp
	allowedEvents map[cqrs.EventType]struct{}
}

// FinishWorkshift implements [edge.ServerInterface].
func (h *EdgeServer) FinishWorkshift(w http.ResponseWriter, r *http.Request) {
	var cmd commands.FinishWorkshiftCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.app.HandleCommand(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// StartWorkshift implements [edge.ServerInterface].
func (h *EdgeServer) StartWorkshift(w http.ResponseWriter, r *http.Request) {
	var cmd commands.StartWorkshiftCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.app.HandleCommand(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Events implements [hardware.ServerInterface].
func (h *EdgeServer) Events(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	stream, err := h.app.EventStream(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	for event := range stream {
		if _, ok := h.allowedEvents[event.Type()]; !ok {
			continue
		}
		data, err := json.Marshal(event)
		if err != nil {
			return
		}
		fmt.Fprintf(w, "event: %s\n", event.Type().String())
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}
}

// FinishDowntime implements [hardware.ServerInterface].
func (h *EdgeServer) FinishDowntime(w http.ResponseWriter, r *http.Request) {
	var cmd commands.FinishDowntimeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.app.HandleCommand(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// StartDowntime implements [hardware.ServerInterface].
func (h *EdgeServer) StartDowntime(w http.ResponseWriter, r *http.Request) {
	var cmd commands.StartDowntimeCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.app.HandleCommand(r.Context(), cmd); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func NewEdgeServer(addr string, app EdgeApp) *EdgeServer {
	srv := &EdgeServer{Server: &http.Server{Addr: addr}, app: app,
		allowedEvents: map[cqrs.EventType]struct{}{
			domain.EventTypeWorkshiftStarted:      {},
			domain.EventTypeWorkshiftFinished:     {},
			domain.EventTypeDowntimeStarted:       {},
			domain.EventTypeDowntimeFinished:      {},
			domain.EventTypeEquipmentLineExpanded: {},
		},
	}
	mux := http.NewServeMux()
	srv.Handler = edge.HandlerWithOptions(srv, edge.StdHTTPServerOptions{
		BaseRouter: mux,
		BaseURL:    "/api",
	})
	return srv
}
