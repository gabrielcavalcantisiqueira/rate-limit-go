package limiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gabrielcavalcantisiqueira/rate-limit-go/internal/config"
)

func TestLimiter_IPBasedLimit(t *testing.T) {
	store := NewInMemoryStore()

	defaultCfg := config.RateLimitConfig{
		Limit:    2,
		Window:   1 * time.Second,
		LockTime: 5 * time.Second,
	}
	tokenCfgs := map[string]config.RateLimitConfig{}

	limiter := NewRateLimiter(store, defaultCfg, tokenCfgs)

	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "1.2.3.4:1234" // simula IP

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Errorf("Expected 200, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Errorf("Expected 200, got %d", rr.Code)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 429 {
		t.Errorf("Expected 429 Too Many Requests, got %d", rr.Code)
	}
}

func TestLimiter_TokenOverridesIP(t *testing.T) {
	store := NewInMemoryStore()

	defaultCfg := config.RateLimitConfig{
		Limit:    1,
		Window:   1 * time.Second,
		LockTime: 1 * time.Second,
	}

	tokenCfgs := map[string]config.RateLimitConfig{
		"abc123": {
			Limit:    5,
			Window:   1 * time.Second,
			LockTime: 2 * time.Second,
		},
	}

	limiter := NewRateLimiter(store, defaultCfg, tokenCfgs)

	handler := limiter.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "5.6.7.8:1234"
	req.Header.Set("API_KEY", "abc123")

	// 5 requisições com token devem passar
	for i := 0; i < 5; i++ {
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != 200 {
			t.Errorf("Expected 200 for request %d, got %d", i+1, rr.Code)
		}
	}

	// 6ª deve ser bloqueada
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != 429 {
		t.Errorf("Expected 429 Too Many Requests on token limit exceeded, got %d", rr.Code)
	}
}
