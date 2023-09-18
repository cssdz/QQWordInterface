package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/models"
	"net/http"
)

// RecordWord 记录单词
func RecordWord(c *gin.Context) {
	// 绑定请求
	var jsRecordWord models.RecordWordRequest
	err := c.ShouldBindJSON(&jsRecordWord)
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
	if jsRecordWord.ID == "" {
		zap.L().Error("未给出用户id")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "必须给出用户id",
		})
		return
	}

	// 检查用户是否存在
	err = mysql.CheckUser(jsRecordWord.ID)
	if err != nil {
		zap.L().Info("该用户未注册")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "该用户未注册",
		})
		return
	}

	// 判断单词类型
	if jsRecordWord.Type == "record" {
		// 判断单词是否在字典里
		no, err := mysql.GetNoByWord(jsRecordWord.Word)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "词典中暂无该单词，请期待后续功能完善",
			})
			return
		}

		// 判断记录是否已经存在
		err = mysql.GetNoInUserRecord(no)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "单词记录已存在",
			})
			return
		}

		// 将记录导入用户
		err = mysql.InsertUserRecord(jsRecordWord.ID, no)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "单词导入失败",
			})
			return
		}

		err = mysql.UpdateUser(jsRecordWord.ID, 1, 0)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "单词导入失败",
			})
			return
		}

		zap.L().Info("单词导入成功", zap.String("id", jsRecordWord.ID), zap.String("word", jsRecordWord.Word))
		c.JSON(http.StatusOK, gin.H{
			"status": 1,
			"msg":    "单词导入成功",
		})
	}
}
