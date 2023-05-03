// Server rest api.
package rest

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/acme/autocert"
)

// Route.
type router interface {
	Run(string) error
}

type server struct {
	router *chi.Mux
	config *core.Config
}

// Create new route.
func New(conf *core.Config, h *chi.Mux) (router, error) {
	return &server{
		router: h,
		config: conf,
	}, nil
}

// Function server Run.
func (s *server) Run(addresPort string) error {
	log.Println("START SERVER ", addresPort, s.config.ServerTimeout)
	server := &http.Server{
		Addr:              addresPort,
		ReadHeaderTimeout: time.Duration(s.config.ServerTimeout) * time.Second,
		Handler:           s.router,
	}
	if s.config.EnableHTTPS != "" {
		certManager := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist("ushorten.ru"),
			Cache:      autocert.DirCache("./certs"),
		}
		server.TLSConfig = &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12,
		}
		go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		// go http.ListenAndServeTLS(":443", "./certs/fullchain.pem", "./certs/privkey.pem", nil)
		return server.ListenAndServeTLS("", "")
	}
	return server.ListenAndServe()
}
