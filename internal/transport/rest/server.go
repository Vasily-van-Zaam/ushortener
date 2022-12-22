package rest

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/go-chi/chi/v5"
)

type Router interface {
	Run(string) error
}

type Server struct {
	router *chi.Mux
	config *core.Config
}

func NewServer(h *handler.Handlers, conf *core.Config) (Router, error) {
	r := chi.NewRouter()
	h.InitAPI(r)
	return &Server{
		router: r,
		config: conf,
	}, nil
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func gzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("MIDLEWARE", next)
		// проверяем, что клиент поддерживает gzip-сжатие

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		// создаём gzip.Writer поверх текущего w
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}

		next.ServeHTTP(gzr, r)
	})
}
func (s *Server) Run(addresPort string) error {
	log.Println("START SERVER ", addresPort, s.config.ServerTimeout)
	server := &http.Server{
		Addr:              addresPort,
		ReadHeaderTimeout: time.Duration(s.config.ServerTimeout) * time.Second,
		Handler:           gzipHandle(s.router),
	}
	return server.ListenAndServe()
}
