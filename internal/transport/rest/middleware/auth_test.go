package middleware

import (
	"crypto/cipher"
	"encoding/hex"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
)

func TestAuth_Handle(t *testing.T) {
	type fields struct {
		Config *core.Config
	}
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				Config: tt.fields.Config,
			}
			if got := a.Handle(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCripto_Encript(t *testing.T) {
	type fields struct {
		nonce  *[]byte
		aesgcm *cipher.AEAD
	}
	type args struct {
		src []byte
	}
	key := "secret0"
	cripto, e := NewCripto([]byte(key))
	log.Println(cripto, e)
	dst := cripto.Encript([]byte("hello12323212"))
	log.Println(hex.EncodeToString(dst))
	tex := hex.EncodeToString(dst)
	xet, _ := hex.DecodeString(tex)
	src, err := cripto.Dencript(xet)

	log.Println(string(src), err, []byte("hello12323212"))
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
	}{
		///
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Cripto{
				nonce:  tt.fields.nonce,
				aesgcm: tt.fields.aesgcm,
			}
			if got := e.Encript(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cripto.Encript() = %v, want %v", got, tt.want)
			}
		})
	}
}
