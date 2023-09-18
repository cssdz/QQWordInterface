package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/models"
	"net/http"
)

// Register 注册/设置
func Register(c *gin.Context) {
	// 绑定请求
	var jsRegister models.Register
	err := c.ShouldBindJSON(&jsRegister)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON(&p)", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "系统错误",
		})
		return
	}

	// 判断信息内容
	// id不能为空
	if jsRegister.ID == "" {
		zap.L().Error("未给出用户id")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "必须给出用户id",
		})
		return
	}

	// 名字和个数不能同时为空
	if jsRegister.Num == -1 && jsRegister.Name == "" {
		zap.L().Error("名字和单词个数不能同时为空")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "名字和单词个数不能同时为空",
		})
		return
	}

	// 检查用户是否注册
	err = mysql.CheckUser(jsRegister.ID)
	if err != nil {
		zap.L().Info("该用户未被注册")
		// 若有一个字段未给值
		if jsRegister.Num == -1 || jsRegister.Name == "" {
			zap.L().Error("用户信息不完整", zap.String("id", jsRegister.ID))
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "用户信息不完整",
			})
			return
		} else { //注册用户
			err = mysql.RegisterUser(jsRegister)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
				return
			}
		}
	} else { // 修改用户信息
		err = mysql.UpdateUserInfo(jsRegister)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
	}

	zap.L().Info("注册/设置成功", zap.String("id", jsRegister.ID))
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "注册/设置成功",
	})
	return
}
