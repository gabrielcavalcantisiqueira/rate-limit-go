package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gabrielcavalcantisiqueira/rate-limit-go/internal/config"
	"github.com/gabrielcavalcantisiqueira/rate-limit-go/internal/limiter"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load .env: %v", err)
	}
	var store limiter.Store
	if config.RedisAddr != "" {
		store = limiter.NewRedisStore(config.RedisAddr)
		fmt.Println("Using Redis as store")
	} else {
		store = limiter.NewInMemoryStore()
		fmt.Println("Using in-memory store")
	}
	rateLimiter := limiter.NewRateLimiter(store, config.DefaultConfig, config.TokenConfigs)
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", rateLimiter.Middleware(mux))
}
