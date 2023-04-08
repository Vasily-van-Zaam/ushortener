package middleware_test

// import (
// 	"encoding/hex"
// 	"log"
// 	"net/http"
// 	"reflect"
// 	"testing"

// 	"github.com/Vasily-van-Zaam/ushortener/internal/core"
// 	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/middleware"
// )

// func TestAuth_Handle(t *testing.T) {
// 	type fields struct {
// 		Config *core.Config
// 	}
// 	type args struct {
// 		next http.Handler
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   http.Handler
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			a := &middleware.Auth{
// 				Config: tt.fields.Config,
// 			}
// 			if got := a.Handle(tt.args.next); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Auth.Handle() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestCripto_Encript(t *testing.T) {
// 	type args struct {
// 		key []byte
// 		src []byte
// 	}
// 	key := "secret1"
// 	cripto, e := middleware.NewCripto([]byte(key))
// 	log.Println(cripto, e)
// 	dst := cripto.Encript([]byte("hello12323212"))
// 	log.Println(hex.EncodeToString(dst))
// 	tex := hex.EncodeToString(dst)
// 	xet, _ := hex.DecodeString(tex)
// 	src, err := cripto.Dencript(xet)

// 	log.Println(string(src), err, []byte("hello12323212"))
// 	tests := []struct {
// 		name string
// 		args args
// 		want []byte
// 	}{
// 		{
// 			name: "cripto",
// 			args: args{
// 				key: []byte("secret1"),
// 				src: src,
// 			},
// 			want: nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			e, _ := middleware.NewCripto(tt.args.key)
// 			if got := e.Encript(tt.args.src); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Cripto.Encript() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
