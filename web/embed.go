package web

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

//go:embed all:dist
var dist embed.FS

func UI() (string, http.Handler) {
	sub, _ := fs.Sub(dist, "dist")
	fileServer := http.FileServer(http.FS(sub))
	return "/oee/", http.StripPrefix("/oee/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// если файл не найден — отдаём index.html (SPA)
		if _, err := sub.Open(strings.TrimPrefix(r.URL.Path, "/")); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	}))
}

func API(target *url.URL) (string, http.Handler) {
	return "/oee/api/", http.StripPrefix("/oee/",
		httputil.NewSingleHostReverseProxy(target),
	)
}
