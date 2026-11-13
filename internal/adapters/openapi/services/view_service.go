package services

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/atmosone/cqrs"
	"github.com/atmosone/oee/internal/queries"
	"github.com/atmosone/oee/pkg/openapi"
)

type ViewService struct {
	logger *slog.Logger
	app    cqrs.QueryHandler
}

func NewViewService(logger *slog.Logger, app cqrs.QueryHandler) *ViewService {
	return &ViewService{logger: logger, app: app}
}

// CurrentInterval implements [openapi.ServerInterface].
func (s *ViewService) CurrentInterval(w http.ResponseWriter, r *http.Request, params openapi.CurrentIntervalParams) {
	res, err := s.app.HandleQuery(r.Context(), queries.CurrentIntervalQuery{
		CurrentIntervalParams: params,
	})
	if err != nil {
		if errors.Is(err, queries.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.logger.Error("failed to read", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// CurrentIntervalsList implements [openapi.ServerInterface].
func (s *ViewService) CurrentIntervalsList(w http.ResponseWriter, r *http.Request) {
	res, err := s.app.HandleQuery(r.Context(), queries.CurrentIntervalsListQuery{})
	if err != nil {
		s.logger.Error("failed to read", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// CurrentWorkshift implements [openapi.ServerInterface].
func (s *ViewService) CurrentWorkshift(w http.ResponseWriter, r *http.Request) {
	res, err := s.app.HandleQuery(r.Context(), queries.CurrentWorkshiftQuery{})
	if err != nil {
		if errors.Is(err, queries.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.logger.Error("failed to read", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// LineTopology implements [openapi.ServerInterface].
func (s *ViewService) LineTopology(w http.ResponseWriter, r *http.Request, params openapi.LineTopologyParams) {
	res, err := s.app.HandleQuery(r.Context(), queries.LineTopologyQuery{
		LineTopologyParams: params,
	})
	if err != nil {
		s.logger.Error("failed to read", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// LinesList implements [openapi.ServerInterface].
func (s *ViewService) LinesList(w http.ResponseWriter, r *http.Request) {
	res, err := s.app.HandleQuery(r.Context(), queries.LinesListQuery{})
	if err != nil {
		s.logger.Error("failed to read", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
