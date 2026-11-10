package openapi

import (
	"net/http"

	"github.com/atmosone/cqrs"
)

func New(addr string, app cqrs.App) *http.Server {
	return &http.Server{Addr: addr}
}
