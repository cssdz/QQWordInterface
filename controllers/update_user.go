package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/models"
	"net/http"
	"strconv"
)

// UpdateUser 更新用户记录
func UpdateUser(c *gin.Context) {
	// 绑定请求
	var jsUpdateUser models.UpdateUserRequest
	err := c.ShouldBindJSON(&jsUpdateUser)
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
	if jsUpdateUser.ID == "" {
		zap.L().Error("未给出用户id")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "必须给出用户id",
		})
		return
	}

	// 检查用户是否注册
	err = mysql.CheckUser(jsUpdateUser.ID)
	if err != nil {
		zap.L().Info("该用户未被注册", zap.String("id", jsUpdateUser.ID))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "该用户未注册",
		})
		return
	}

	// 单词编号必须为数字
	_, err = strconv.Atoi(jsUpdateUser.No)
	if err != nil {
		zap.L().Error("单词编号必须为数字", zap.String("id", jsUpdateUser.ID))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "单词编号必须为数字",
		})
		return
	}

	// status必须为0或1
	status, err := strconv.Atoi(jsUpdateUser.Status)
	if err != nil {
		zap.L().Error("变量格式转换错误", zap.String("id", jsUpdateUser.ID), zap.Int("status", status))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "变量格式转换错误",
		})
		return
	}
	if status != 0 && status != 1 {
		zap.L().Error("单词状态必须为0或1", zap.String("id", jsUpdateUser.ID), zap.Int("status", status))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "单词状态必须为0或1",
		})
		return
	}

	// 查询记录是否存在
	err = mysql.SelUserRecord(jsUpdateUser.ID, jsUpdateUser.No)
	if err != nil { // 如果记录不存在且单词状态为0，则插入一条新记录
		turn := 0
		zap.L().Info("记录不存在", zap.String("id", jsUpdateUser.ID), zap.String("no", jsUpdateUser.No))
		if status == 0 {
			err = mysql.InsertUserRecord(jsUpdateUser.ID, jsUpdateUser.No)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
				return
			}
		}
		// 更新用户数据
		err = mysql.UpdateUser(jsUpdateUser.ID, jsUpdateUser.No, status, turn)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
	} else {
		turn := 1
		zap.L().Info("记录存在", zap.String("id", jsUpdateUser.ID), zap.String("no", jsUpdateUser.No))
		if status == 0 {
			err = mysql.UpdateUserRecord(jsUpdateUser.ID, jsUpdateUser.No)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
				return
			}
		}
		// 更新用户数据
		err = mysql.UpdateUser(jsUpdateUser.ID, jsUpdateUser.Status, status, turn)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}
	}

	zap.L().Info("用户记录更新成功", zap.String("id", jsUpdateUser.ID))
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "用户记录更新成功",
	})
}
