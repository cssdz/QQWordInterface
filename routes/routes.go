package routes

import (
	"go_web/controllers"
	"go_web/logger"
	"net/http"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	if viper.GetString("app.mode") != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, viper.GetString("app.version"))
	})

	// 添加路由前缀
	v1 := r.Group("/api/v1")

	// 将单词书导入数据库
	v1.POST("/import_vocabulary", Check(), controllers.WordBookToMySQL)

	// 注册/设置
	v1.POST("/register", Check(), controllers.Register)

	// 获取用户信息
	v1.POST("/user_info", Check(), controllers.GetUserInfo)

	// 获取单词
	v1.POST("/get_words", Check(), controllers.GetWords)

	// 用户记录更新
	v1.POST("/update_user", Check(), controllers.UpdateUser)

	// 单词记录
	v1.POST("/record_word", Check(), controllers.RecordWord)

	return r
}
