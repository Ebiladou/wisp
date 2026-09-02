package routes

import (
	"time"

	"github.com/Ebiladou/wisp/internal/middleware"
)

var (
	defaultRateLimit = middleware.RateLimitConfig{
		Limit:  30,
		Window: time.Minute,
	}

	authRateLimit = middleware.RateLimitConfig{
		Limit:  10,
		Window: time.Minute,
	}

	sensitiveAuthLimit = middleware.RateLimitConfig{
		Limit:  3,
		Window: time.Minute,
	}
)
