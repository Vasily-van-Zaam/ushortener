package handler

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// /// mock.
type ServiceMock struct {
	core.AUTHService
}

func (s *ServiceMock) GetURL(ctx context.Context, id string) (string, error) {
	switch id {
	case "1":
		{
			return "http://example.com/link1", nil
		}
	case "2":
		{
			return "http://example.com/link2", nil
		}
	default:
		{
			return "", errors.New("not Found")
		}
	}
}
func (s *ServiceMock) SetURL(ctx context.Context, link string) (string, error) {
	switch link {
	case "http://example.com/link1":
		{
			return "http://localhost:8080/" + "1", nil
		}
	case "http://example.com/link2":
		{
			return "http://localhost:8080/" + "2", nil
		}
	default:
		{
			return "", nil
		}
	}
}

func (s *ServiceMock) APISetShorten(
	ctx context.Context, request *core.RequestAPIShorten) (*core.ResponseAPIShorten, error) {
	return &core.ResponseAPIShorten{}, nil
}

func TestShortenerHandler_GetSetURL(t *testing.T) {
	service := ServiceMock{}
	type fields struct {
		service BasicService
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		// TODO: Add test cases.
		{
			name: "set short link 1: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					bytes.NewReader([]byte("http://example.com/link1")),
					// strings.NewReader("http://example.com/link1"),
				),
			},
			want: want{
				code:        201,
				response:    `http://localhost:8080/1`,
				contentType: "text/plain",
			},
		},
		{
			name: "set short link 2: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					strings.NewReader("http://example.com/link2"),
				),
			},
			want: want{
				code:        201,
				response:    `http://localhost:8080/2`,
				contentType: "text/plain",
			},
		},
		{
			name: "set short link empty post body: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(
					http.MethodPost,
					"/",
					nil,
				),
			},
			want: want{
				code:        201,
				response:    `http://localhost:8080/`,
				contentType: "text/plain",
			},
		},
		{
			name: "get short link 1: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/1", nil),
			},
			want: want{
				code:        307,
				response:    `http://example.com/link1`,
				contentType: "text/plain",
			},
		},
		{
			name: "get short link 2: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/2", nil),
			},
			want: want{
				code:        307,
				response:    `http://example.com/link2`,
				contentType: "text/plain",
			},
		},
		{
			name: "get error: not Found: ",
			fields: fields{
				service: &service,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/", nil),
			},
			want: want{
				code:        400,
				response:    "not Found\n",
				contentType: "text/plain",
			},
		},
	}
	for _, tt := range tests {
		cfg := core.Config{
			ServerAddress: "127.0.0.1:8080/",
			BaseURL:       "http://localhost:8080/",
		}

		t.Run(tt.name, func(t *testing.T) {
			h := &BasicHandler{
				service: &tt.fields.service,
				config:  &cfg,
			}
			r := chi.NewRouter()
			hs := NewHandlers(h, nil)
			hs.InitAPI(r)
			r.ServeHTTP(tt.args.w, tt.args.r)

			///////// chech response //////////
			res := tt.args.w.Result()

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.response, string(resBody))

			log.Println(tt.name, string(resBody), res.StatusCode)
		})
	}
}

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
func main1() {
	var (
		data  []byte         // слайс случайных байт
		hash1 []byte         // хеш с использованием интерфейса hash.Hash
		hash2 [md5.Size]byte // хеш, возвращаемый функцией md5.Sum
	)
	// допишите код
	// 1) генерация data длиной 512 байт
	// 2) вычисление hash1 с использованием md5.New
	// 3) вычисление hash2 функцией md5.Sum

	data = []byte("Hello!! 1234") // make([]byte, 512)
	// _, err := rand.Read(data)
	// if err != nil {
	// 	log.Println(err)
	// }
	h := md5.New()
	h.Write(data)
	hash1 = h.Sum(nil)
	hash2 = md5.Sum(data)
	log.Println(hash2[:4], hash2[4:])
	log.Println(hash2, hash1)
	log.Println(data)
	// ...

	// hash2[:] приводит массив байт к слайсу
	if bytes.Equal(hash1, hash2[:]) {
		fmt.Println("Всё правильно! Хеши равны")
	} else {
		fmt.Println("Что-то пошло не так")
	}
}

var secretkey = []byte("secret key0")

func main2() {
	var (
		data []byte // декодированное сообщение с подписью
		id   uint32 // значение идентификатора
		err  error
		sign []byte // HMAC-подпись от идентификатора
	)
	msg := "048ff4ea240a9fdeac8f1422733e9f3b8b0291c969652225e25c5f0f9f8da654139c9e21"

	// допишите код
	// 1) декодируйте msg в data
	// 2) получите идентификатор из первых четырёх байт,
	//    используйте функцию binary.BigEndian.Uint32
	// 3) вычислите HMAC-подпись sign для этих четырёх байт

	// ...
	data, err = hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}
	id = binary.BigEndian.Uint32(data[:4])
	h := hmac.New(sha256.New, secretkey)
	h.Write(data[:4])
	sign = h.Sum(nil)

	if hmac.Equal(sign, data[4:]) {
		fmt.Println("Подпись подлинная. ID:", id)
	} else {
		fmt.Println("Подпись неверна. Где-то ошибка", id)
	}
}
func TestCripto(t *testing.T) {
	// log.Println(RandBytes((100)))
	// main1()
	// main2()
	// log.Println(aes.BlockSize)
	const (
		password = "x35k9f"
		msg      = `0ba7cd8c624345451df4710b81d1a349ce401e61bc7eb704ca` +
			`a84a8cde9f9959699f75d0d1075d676f1fe2eb475cf81f62ef` +
			`f701fee6a433cfd289d231440cf549e40b6c13d8843197a95f` +
			`8639911b7ed39a3aec4dfa9d286095c705e1a825b10a9104c6` +
			`be55d1079e6c6167118ac91318fe`
	)
	key := sha256.Sum256([]byte(password))
	log.Println(key)
	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	encrypted, err := hex.DecodeString(msg)
	if err != nil {
		panic(err)
	}

	decrypted, err := aesgcm.Open(nil, nonce, encrypted, nil)
	if err != nil {
		panic(err)
	}
	log.Print(string(decrypted))
	// src := []byte("Этюд в розовых тонах") // данные, которые хотим зашифровать
	// fmt.Printf("original: %s\n", src)

	// // будем использовать AES256, создав ключ длиной 32 байта
	// key, err := generateRandom(2 * aes.BlockSize) // ключ шифрования
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return
	// }

	// aesblock, err := aes.NewCipher(key)
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return
	// }

	// aesgcm, err := cipher.NewGCM(aesblock)
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return
	// }

	// // создаём вектор инициализации
	// nonce, err := generateRandom(aesgcm.NonceSize())
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return
	// }

	// dst := aesgcm.Seal(nil, nonce, src, nil) // зашифровываем
	// fmt.Printf("encrypted: %x\n", dst)

	// src2, err := aesgcm.Open(nil, nonce, dst, nil) // расшифровываем
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return
	// }
	// fmt.Printf("decrypted: %s\n", src2)
}
