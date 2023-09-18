package routes

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// Check 检验token中间件
func Check() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		// 是否为空
		if authHeader == "" {
			zap.L().Error("未找到token", zap.String("authHeader", authHeader))
			c.JSON(http.StatusOK, gin.H{
				"status": -7,
				"msg":    "未找到token",
			})
			c.Abort()
			return
		}

		// 将uid字符串转换为数字
		token := authHeader

		TrueToken := "cs7949335038094359bb67a16f1e1808fdbde8e256sdz"

		// 判断token是否相等
		if token != TrueToken {
			zap.L().Error("无效的Token", zap.String("token", token), zap.String("TrueToken", TrueToken))
			c.JSON(http.StatusOK, gin.H{
				"status": -7,
				"msg":    "无效的Token",
			})
			c.Abort()
			return
		}

		zap.L().Info("Authorization Accept", zap.String("header", authHeader))
		return
	}
}

// CORSMiddleware 跨域中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
