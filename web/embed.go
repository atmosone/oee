package web

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

//go:embed all:dist
var dist embed.FS

func UI() (string, http.Handler) {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	return "/oee/", http.StripPrefix("/oee/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Очищаем путь (убираем начальный слеш, если есть)
		p := strings.TrimPrefix(r.URL.Path, "/")

		// 2. Если пусто (запрос на /oee/) — ставим index.html
		if p == "" {
			p = "index.html"
		}

		// 3. Пытаемся открыть файл
		f, err := sub.Open(p)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				// SPA fallback: файла нет — отдаем index.html
				p = "index.html"
				f, err = sub.Open(p)
				if err != nil {
					http.Error(w, "index.html not found in dist", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		defer f.Close()

		// 4. Проверяем, что это не директория
		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		// 5. Определяем MIME-тип по расширению
		ext := path.Ext(p)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		w.Header().Set("Content-Type", contentType)

		// 6. Отдаем файл (используем http.ServeContent для поддержки Range-запросов)
		http.ServeContent(w, r, p, stat.ModTime(), f.(io.ReadSeeker))
	}))
}
func API(target *url.URL) (string, http.Handler) {
	return "/oee/api/", http.StripPrefix("/oee/",
		httputil.NewSingleHostReverseProxy(target),
	)
}
