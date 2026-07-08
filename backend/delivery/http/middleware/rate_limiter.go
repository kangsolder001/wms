package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pkgredis "wms/pkg/redis"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	redis  pkgredis.RedisClient
	limit  int
	window time.Duration
}

func NewRateLimiter(redisClient pkgredis.RedisClient, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		redis:  redisClient,
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) isRateLimited(ctx context.Context, key string) (bool, error) {
	current, err := rl.redis.Incr(ctx, key)
	if err != nil {
		return false, err
	}

	if current == 1 {
		if err := rl.redis.Expire(ctx, key, rl.window); err != nil {
			return false, err
		}
	}

	return current > int64(rl.limit), nil
}

func (rl *RateLimiter) LoginLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("login_rate:%s", ip)

		ctx := context.Background()

		limited, err := rl.isRateLimited(ctx, key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
			return
		}

		if limited {
			ttl, _ := rl.redis.TTL(ctx, key)
			seconds := int(ttl.Seconds())
			if seconds <= 0 {
				seconds = 60
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("too many login attempts, please try again in %d seconds", seconds),
			})
			return
		}

		c.Next()
	}
}
