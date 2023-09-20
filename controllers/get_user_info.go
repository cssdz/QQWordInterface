package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/models"
	"net/http"
	"time"
)

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	// 绑定请求
	var jsUserInfoRequest models.UserInfoRequest
	err := c.ShouldBindJSON(&jsUserInfoRequest)
	if err != nil {
		zap.L().Error("c.ShouldBindJSON(&p)", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "系统错误",
		})
		return
	}

	// 判断该用户是否注册
	err = mysql.CheckUser(jsUserInfoRequest.ID)
	if err != nil {
		zap.L().Warn("该用户未注册", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "该用户未注册",
			"bool":   0,
		})
		return
	}

	// 获取用户数据
	var userInfo models.GetUserInfo
	userInfo, err = mysql.GetUserInfo(jsUserInfoRequest.ID)
	if err != nil {
		zap.L().Error("用户数据获取失败", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	// 修改数据格式
	// 修改go从数据库获取的时间格式
	strTime := userInfo.RegisterTime
	t1, _ := time.Parse(time.RFC3339, strTime)
	s1 := t1.Format("2006-01-02 15:04:05")
	userInfo.RegisterTime = s1

	strTime = userInfo.UpdateTime
	t2, _ := time.Parse(time.RFC3339, strTime)
	s2 := t2.Format("2006-01-02 15:04:05")
	userInfo.UpdateTime = s2

	userInfo.UseTimes = fmt.Sprintf("%s次", userInfo.UseTimes)
	userInfo.WordTimes = fmt.Sprintf("%s次", userInfo.WordTimes)
	userInfo.RecordWordNum = fmt.Sprintf("记录了%s个单词", userInfo.RecordWordNum)
	userInfo.MaxAchievement = fmt.Sprintf("最多连续学习%s天", userInfo.MaxAchievement)
	userInfo.LatestAchievement = fmt.Sprintf("最近连续学习%s天", userInfo.LatestAchievement)

	zap.L().Info("注册/设置成功", zap.String("id", jsUserInfoRequest.ID))
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "注册/设置成功",
		"bool":   1,
		"data":   userInfo,
	})
	return
}
