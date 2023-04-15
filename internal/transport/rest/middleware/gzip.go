package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type zip struct {
	Config *core.Config
}

func NewGzip(conf *core.Config) *zip {
	return &zip{
		Config: conf,
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func (g *zip) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			rg, errs := gzip.NewReader(r.Body)
			if errs != nil {
				log.Println("ERROR GZIP", errs)
				w.WriteHeader(http.StatusBadRequest)
				_, err := w.Write([]byte("error request data gzip "))
				if err != nil {
					log.Println("EERROR WRITE ", err)
				}
				next.ServeHTTP(w, r)
				return
			}
			bodyBytes, _ = io.ReadAll(rg)
		} else {
			bodyBytes, _ = io.ReadAll(r.Body)
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			g.Config.LogRequest(w, r, string(bodyBytes))
			return
		}

		w.Header().Set("Content-Encoding", "gzip")

		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		g.Config.LogRequest(w, r, string(bodyBytes))
		next.ServeHTTP(gzr, r)
	})
}
