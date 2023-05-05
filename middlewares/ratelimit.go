package middlewares

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果获取不到令牌就返回限流响应
		if bucket.TakeAvailable(1) == 0 {
			c.JSON(http.StatusTooManyRequests, "请求数量较多，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
