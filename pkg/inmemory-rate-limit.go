package pkg

import (
	"net"
	"net/http"
	"time"
)

const (
	limit    = 5                // Máximo de requisições permitidas
	window   = 30 * time.Second // Intervalo de tempo
	lockTime = 1 * time.Minute  // Tempo para limpar IPs inativos
)

var clientLocked map[string]time.Time = make(map[string]time.Time)
var clientRequests map[string][]time.Time = make(map[string][]time.Time)

func UseRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		now := time.Now()
		if lockedUntil, exists := clientLocked[ip]; exists && now.Before(lockedUntil) {
			http.Error(w, "429 - Too Many Requests", http.StatusTooManyRequests)
			return
		}
		requests, exists := clientRequests[ip]
		if !exists {
			requests = make([]time.Time, 0)
		}
		var requestsWithinTheWindow []time.Time
		for _, t := range requests {
			if now.Sub(t) <= window {
				requestsWithinTheWindow = append(requestsWithinTheWindow, t)
			}
		}
		if len(requestsWithinTheWindow) >= limit {
			clientLocked[ip] = now.Add(lockTime)
			http.Error(w, "429 - Too Many Requests", http.StatusTooManyRequests)
			return
		}
		requestsWithinTheWindow = append(requestsWithinTheWindow, now)
		clientRequests[ip] = requestsWithinTheWindow
		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return fwd
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
