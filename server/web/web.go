package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
)

//go:embed all:qualification
var qualificationFS embed.FS

func handler(dist fs.FS) http.Handler {
	server := http.FileServer(http.FS(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if path == "" {
			path = "."
		}
		f, err := dist.Open(path)
		if err != nil {
			data, err := fs.ReadFile(dist, "index.html")
			if err != nil {
				http.Error(w, "index.html not found", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Write(data)
			return
		}
		f.Close()
		server.ServeHTTP(w, r)
	})
}

func NewQualificationHandler() http.Handler {
	dist, err := fs.Sub(qualificationFS, "qualification")
	if err != nil {
		log.Fatal(err)
	}
	return handler(dist)
}
