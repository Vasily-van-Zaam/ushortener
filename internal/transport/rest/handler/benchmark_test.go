package handler_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/core"
	"github.com/Vasily-van-Zaam/ushortener/internal/transport/rest/handler"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MockBencnmarkBasicService struct {
	data    map[string]string
	mu      sync.RWMutex
	baseUrl string
}

func NewMockBencnmarkBasicService() *MockBencnmarkBasicService {
	data := make(map[string]string)
	return &MockBencnmarkBasicService{
		data:    data,
		mu:      sync.RWMutex{},
		baseUrl: "http://site.ru/",
	}
}
func (s *MockBencnmarkBasicService) CreateUser() {

}

func (m *MockBencnmarkBasicService) GetURL(ctx context.Context, link string) (string, error) {
	m.mu.RLock()
	id := strings.ReplaceAll(link, m.baseUrl, "")
	resp := m.data[id]
	m.mu.RUnlock()
	if resp == "" {
		return "", errors.New("not Found")
	}
	return resp, nil
}
func (m *MockBencnmarkBasicService) SetURL(ctx context.Context, link string) (string, error) {
	m.mu.Lock()
	id := uuid.New().String()
	m.data[id] = link
	m.mu.Unlock()
	return m.baseUrl + id, nil
}
func (m *MockBencnmarkBasicService) Ping(ctx context.Context) error {
	return nil
}

func BenchmarkSetAdd1000Urls(b *testing.B) {
	countURL := 1000
	userDomain := "https://google.com/"
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	basicService := NewMockBencnmarkBasicService()
	conf := &core.Config{
		ServerAddress: "http://localhost:8080",
		BaseURL:       basicService.baseUrl,
	}
	h := handler.NewBasic(basicService, conf)

	newURL := userDomain + uuid.New().String() + "/" + uuid.New().String()

	urls := make([]string, countURL)
	for i := range urls {
		urls[i] = userDomain + uuid.New().String() + "/" + uuid.New().String()
	}
	r := chi.NewRouter()
	hs := handler.NewHandlers(h, nil)
	hs.InitAPI(r)

	b.ResetTimer()

	for _ = range urls {
		b.StopTimer()
		a := args{
			w: httptest.NewRecorder(),
			r: httptest.NewRequest(
				http.MethodPost,
				"/",
				strings.NewReader(newURL),
			),
		}
		b.StartTimer()
		r.ServeHTTP(a.w, a.r)
	}
}
