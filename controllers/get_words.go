package controllers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/logic"
	"go_web/models"
	"net/http"
)

// GetWords 获取单词
func GetWords(c *gin.Context) {
	// 绑定请求
	var jsGetWordsRequest models.GetWordsRequest
	err := c.ShouldBindJSON(&jsGetWordsRequest)
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
	if jsGetWordsRequest.ID == "" {
		zap.L().Error("未给出用户id")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "必须给出用户id",
		})
		return
	}

	// 检查用户是否存在
	err = mysql.CheckUser(jsGetWordsRequest.ID)
	if err != nil {
		zap.L().Info("该用户未注册")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "该用户未注册",
		})
		return
	}

	// 判断获取的单词数量是否大于0
	if jsGetWordsRequest.Num < 0 {
		zap.L().Info("获取的单词数不合法")
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "获取的单词数不合法",
		})
		return
	} else if jsGetWordsRequest.Num == 0 {
		// 查询用户默认的单词数
		jsGetWordsRequest.Num, err = mysql.GetUserNum(jsGetWordsRequest.ID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "获取单词失败",
			})
			return
		}
		zap.L().Info("获取用户默认单词数成功", zap.Int("num", jsGetWordsRequest.Num))
	}

	logic.UpdateAchievement(jsGetWordsRequest.ID)
	var WordString []models.WordString
	// 判断获取单词的类型
	if jsGetWordsRequest.Type == "record" {
		// 查询复习单词的数量
		recordNum, err := mysql.GetRecordWordsNum(jsGetWordsRequest.ID)
		if err != nil {
			zap.L().Error("获取复习单词数量失败", zap.String("id", jsGetWordsRequest.ID))
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "获取的单词数不合法",
			})
			return
		}

		// 如果记录不会的单词数为0，则直接从单词书中拿新单词
		if recordNum == 0 {
			// 查询该用户要背的下一个单词
			var words []models.Word
			nextNo, err := mysql.GetNextWordNo(jsGetWordsRequest.ID)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
				return
			}

			// 获取新单词
			maxNo := nextNo + jsGetWordsRequest.Num - 1
			words, err = mysql.GetNewWords(jsGetWordsRequest.ID, nextNo, maxNo)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"status": -1,
					"msg":    err.Error(),
				})
				return
			}

			WordString = logic.ToString(words)
		} else { // 如果记录的单词数小于请求的单词数，则先背记录的单词再背新单词
			var recordWords []models.Word
			num := jsGetWordsRequest.Num - recordNum
			if num > 0 {
				// 获取记录的单词
				recordWords, err = mysql.GetRecordWords(jsGetWordsRequest.ID, recordNum)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"status": -1,
						"msg":    err.Error(),
					})
					return
				}

				// 查询该用户要背的下一个单词
				nextNo, err := mysql.GetNextWordNo(jsGetWordsRequest.ID)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"status": -1,
						"msg":    err.Error(),
					})
					return
				}

				// 获取新单词
				var newWords []models.Word
				maxNo := nextNo + num - 1
				newWords, err = mysql.GetNewWords(jsGetWordsRequest.ID, nextNo, maxNo)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"status": -1,
						"msg":    err.Error(),
					})
					return
				}

				words := logic.Combine(recordWords, newWords)
				WordString = logic.ToString(words)
			} else { // 如果记录的单词数大于等于请求数
				// 获取记录的单词
				recordWords, err = mysql.GetRecordWords(jsGetWordsRequest.ID, recordNum)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"status": -1,
						"msg":    err.Error(),
					})
					return
				}

				WordString = logic.ToString(recordWords)
			}
		}
	} else if jsGetWordsRequest.Type == "new" {
		// 查询该用户要背的下一个单词
		var words []models.Word
		nextNo, err := mysql.GetNextWordNo(jsGetWordsRequest.ID)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}

		// 获取新单词
		maxNo := nextNo + jsGetWordsRequest.Num - 1
		words, err = mysql.GetNewWords(jsGetWordsRequest.ID, nextNo, maxNo)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    err.Error(),
			})
			return
		}

		WordString = logic.ToString(words)
	} else {
		zap.L().Error("输入的类型有误", zap.String("id", jsGetWordsRequest.ID), zap.String("type", jsGetWordsRequest.Type))
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "输入的类型有误",
		})
		return
	}

	// 获取成功使用户的请求记录+1同时更新操作时间
	err = mysql.UseTimesAddOne(jsGetWordsRequest.ID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		return
	}

	zap.L().Info("单词获取成功", zap.String("id", jsGetWordsRequest.ID))
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "单词获取成功",
		"data":   WordString,
	})
	return
}
