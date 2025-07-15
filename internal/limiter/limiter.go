package limiter

import (
	"net/http"
	"strings"
	"time"

	"github.com/gabrielcavalcantisiqueira/rate-limit-go/internal/config"
)

type RateLimiter struct {
	store        Store
	defaultConf  config.RateLimitConfig
	tokenConfigs map[string]config.RateLimitConfig
}

func NewRateLimiter(store Store, def config.RateLimitConfig, tokens map[string]config.RateLimitConfig) *RateLimiter {
	return &RateLimiter{
		store:        store,
		defaultConf:  def,
		tokenConfigs: tokens,
	}
}

func (r *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		id, cfg := r.extractIDAndConfig(req)
		if r.store.IsLocked(id) {
			http.Error(w, "429 - Too Many Requests: you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		if !r.store.Allow(id, cfg.Limit, cfg.Window) {
			r.store.SetLock(id, time.Now().Add(cfg.LockTime))
			http.Error(w, "429 - Too Many Requests: you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (r *RateLimiter) extractIDAndConfig(req *http.Request) (string, config.RateLimitConfig) {
	token := strings.TrimSpace(req.Header.Get("API_KEY"))
	if token != "" {
		if cfg, ok := r.tokenConfigs[token]; ok {
			return token, cfg
		}
	}
	ip := getClientIP(req)
	return ip, r.defaultConf
}

func getClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		parts := strings.Split(fwd, ",")
		return strings.TrimSpace(parts[0])
	}
	ip := r.RemoteAddr
	if i := strings.LastIndex(ip, ":"); i != -1 {
		return ip[:i]
	}
	return ip
}
