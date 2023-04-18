// Auth middleware
package middleware

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/http"
	"time"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/google/uuid"
)

// Implementation.
type Auth interface {
	Handle(next http.Handler) http.Handler
}

// Struture.
type auth struct {
	Config  *core.Config
	service core.AUTHService
	cripto  cripter
}

// Create a new AUTH middleware.
func NewAuth(conf *core.Config, service core.AUTHService) Auth {
	c, err := NewCripto([]byte(conf.SecretKey))
	if err != nil {
		log.Fatal(err)
	}
	return &auth{
		Config:  conf,
		service: service,
		cripto:  c,
	}
}

// Handle response, request.
func (a *auth) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newReq := a.generatCookie(w, r)
		if next != nil {
			next.ServeHTTP(w, newReq)
		}
	})
}

// Generate new cookies.
func (a *auth) generatCookie(w http.ResponseWriter, r *http.Request) *http.Request {
	cookie := r.Cookies()
	id := uuid.New().String()
	criptID := a.cripto.Encript([]byte(id))

	newCookie := &http.Cookie{
		Path:    "/",
		Name:    string(core.USERDATA),
		Value:   hex.EncodeToString(criptID),
		Expires: time.Now().Add(time.Hour * 24 * time.Duration(a.Config.ExpiresDayCookie)),
	}
	if len(cookie) == 0 {
		http.SetCookie(w, newCookie)
		return r.WithContext(setContext(r, core.User{ID: id}))
	}
	for _, c := range cookie {
		// return r.WithContext(setContext(r, core.User{ID: "3f295f52-fecf-43d0-b5dc-856b817ab820"}))
		if c.Name == string(core.USERDATA) {
			var v []byte
			byteID, err1 := hex.DecodeString(c.Value)
			v, err := a.cripto.Dencript(byteID)
			if err != nil {
				log.Println("ERROR cripto, sended new cookie", err1, err)
				http.SetCookie(w, newCookie)
				return r.WithContext(setContext(r, core.User{ID: id}))
			}
			return r.WithContext(setContext(r, core.User{ID: string(v)}))
		}
	}

	return r.WithContext(setContext(r, core.User{ID: hex.EncodeToString(criptID)}))
}

// Set conext.
func setContext(r *http.Request, user core.User) context.Context {
	return context.WithValue(
		r.Context(),
		core.USERDATA,
		user,
	)
}

// Cripto.
type cripter interface {
	Encript(src []byte) []byte
	Dencript(dst []byte) ([]byte, error)
}
type cripto struct {
	nonce  *[]byte
	aesgcm *cipher.AEAD
}

// New cripto.
func NewCripto(key []byte) (cripter, error) {
	h := sha256.New()
	h.Write(key)
	k := h.Sum(nil)
	aesblock, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	nonce, err := GenerateRandom(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}
	return &cripto{nonce: &nonce, aesgcm: &aesgcm}, nil
}

// Encript.
func (e *cripto) Encript(src []byte) []byte {
	dst := (*e.aesgcm).Seal(nil, *e.nonce, src, nil)
	return dst
}

// Decript.
func (e *cripto) Dencript(dst []byte) ([]byte, error) {
	res, err := (*e.aesgcm).Open(nil, *e.nonce, dst, nil) // расшифровываем
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Generate Random.
func GenerateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
