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
	*services.ViewService
}

func New(
	router *http.ServeMux,
	management *services.ManagementService,
	edge *services.EdgeService,
	event *services.EventService,
	view *services.ViewService,
) http.Handler {
	return openapi.HandlerWithOptions(
		&server{
			ManagementService: management,
			EdgeService:       edge,
			EventService:      event,
			ViewService:       view,
		},
		openapi.StdHTTPServerOptions{
			BaseRouter: router,
			BaseURL:    "/api",
		})
}
