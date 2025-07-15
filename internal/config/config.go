package config

import (
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type RateLimitConfig struct {
	Limit    int
	Window   time.Duration
	LockTime time.Duration
}

var (
	RedisAddr     string
	DefaultConfig RateLimitConfig
	TokenConfigs  map[string]RateLimitConfig
)

func Load() error {
	viper.SetConfigFile(".env")
	viper.SetDefault("REDIS_ADDR", "")
	viper.SetDefault("DEFAULT_LIMIT", 5)
	viper.SetDefault("DEFAULT_WINDOW", "1s")
	viper.SetDefault("DEFAULT_LOCK", "5m")
	viper.SetDefault("TOKENS_LIMITS", "")
	if err := viper.ReadInConfig(); err != nil {
		_ = err
	}
	viper.AutomaticEnv()
	RedisAddr = viper.GetString("REDIS_ADDR")
	DefaultConfig = RateLimitConfig{
		Limit:    viper.GetInt("DEFAULT_LIMIT"),
		Window:   mustDuration(viper.GetString("DEFAULT_WINDOW")),
		LockTime: mustDuration(viper.GetString("DEFAULT_LOCK")),
	}
	TokenConfigs = make(map[string]RateLimitConfig)
	raw := viper.GetString("TOKENS_LIMITS")
	if raw != "" {
		tokens := strings.Split(raw, ",")
		for _, token := range tokens {
			parts := strings.Split(token, ":")
			if len(parts) != 4 {
				continue
			}
			limit := mustInt(parts[1])
			window := mustDuration(parts[2])
			lock := mustDuration(parts[3])
			TokenConfigs[parts[0]] = RateLimitConfig{
				Limit:    limit,
				Window:   window,
				LockTime: lock,
			}
		}
	}
	return nil
}

func mustInt(val string) int {
	i, _ := strconv.Atoi(val)
	return i
}

func mustDuration(val string) time.Duration {
	d, _ := time.ParseDuration(val)
	return d
}
