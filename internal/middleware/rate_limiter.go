package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"

	"github.com/Ebiladou/wisp/internal/authentication"
	"github.com/Ebiladou/wisp/internal/models"
)

type RateLimitConfig struct {
	Limit  int
	Window time.Duration
}

type RateLimiter struct {
	client *goredis.Client
}

func NewRateLimiter(addr string) (*RateLimiter, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &RateLimiter{
		client: client,
	}, nil
}

func (limiter *RateLimiter) Close() error {
	return limiter.client.Close()
}

func (limiter *RateLimiter) Middleware(config RateLimitConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fmt.Sprintf(
			"rate_limit:ip:%s:%s",
			c.ClientIP(),
			c.FullPath(),
		)

		if user, exists := c.Get(authentication.AuthenticatedUserKey); exists {
			authenticatedUser, ok := user.(*models.User)

			if ok {
				key = fmt.Sprintf(
					"rate_limit:user:%s:%s",
					authenticatedUser.ID,
					c.FullPath(),
				)
			}
		}

		count, err := limiter.client.Incr(
			c.Request.Context(),
			key,
		).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter unavailable",
			})
			c.Abort()
			return
		}

		if count == 1 {
			err := limiter.client.Expire(
				c.Request.Context(),
				key,
				config.Window,
			).Err()

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "rate limiter unavailable",
				})
				c.Abort()
				return
			}
		}

		ttl, err := limiter.client.TTL(
			c.Request.Context(),
			key,
		).Result()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "rate limiter unavailable",
			})
			c.Abort()
			return
		}

		if count > int64(config.Limit) {
			if ttl > 0 {
				c.Header(
					"Retry-After",
					fmt.Sprintf("%d", int(ttl.Seconds())),
				)
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
