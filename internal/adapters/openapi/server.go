package openapi

import (
	"net/http"

	"github.com/atmosone/oee/internal/adapters/openapi/services"
	"github.com/atmosone/oee/pkg/openapi"
)

type server struct {
	*services.ManagementService
	*services.EdgeService
	*services.EventService
}

func New(
	router *http.ServeMux,
	management *services.ManagementService,
	edge *services.EdgeService,
	event *services.EventService,
) http.Handler {
	return openapi.HandlerWithOptions(
		&server{ManagementService: management, EdgeService: edge, EventService: event},
		openapi.StdHTTPServerOptions{
			BaseRouter: router,
			BaseURL:    "/api",
		})
}
