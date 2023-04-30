package middlewares

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带token有三种方式 1、放在请求头 2、放在请求体 3、放在URI
		// 此处假定放在Header的Authorization中， 并使用Bearer开头
		// 这里具体实现方式要依据实际业务来定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		if err := mc.Valid(); err != nil {
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将请求的user_id保存到请求的上下文中
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() // 后续请求处理函数可以用c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}
