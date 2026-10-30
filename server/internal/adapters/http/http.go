package http

import "net/http"

type EdgeServer struct{ *http.Server }

func NewEdgeServer(addr string, api http.Handler) *EdgeServer {
	return &EdgeServer{Server: &http.Server{Addr: addr, Handler: api}}
}

type CmsServer struct{ *http.Server }

func NewCmsServer(addr string, api http.Handler) *CmsServer {
	return &CmsServer{Server: &http.Server{Addr: addr, Handler: api}}
}
