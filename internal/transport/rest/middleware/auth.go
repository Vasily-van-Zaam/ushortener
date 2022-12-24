package middleware

import (
	"net/http"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

type Auth struct {
	Config *core.Config
}

func NewAuth(conf *core.Config) *Auth {
	return &Auth{
		Config: conf,
	}
}
func (a *Auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
