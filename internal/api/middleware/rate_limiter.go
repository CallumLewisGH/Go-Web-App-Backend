package middleware

import (
	"os"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.JSON(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

func NewRateLimiter(requestsPer uint, timeUnit time.Duration, redisURL string) gin.HandlerFunc {
	if redisURL == "" {
		// Fall back to env config if no URL provided
		_ = godotenv.Load("/home/callum/Desktop/Go-Web-App-Backend/.dev.env") //Not configured yet
		redisURL = os.Getenv("REDIS_URL")
	}

	store := ratelimit.RedisStore(&ratelimit.RedisOptions{
		RedisClient: redis.NewClient(&redis.Options{
			Addr: redisURL,
		}),
		Rate:  timeUnit,
		Limit: requestsPer,
	})

	mw := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	return mw
}
