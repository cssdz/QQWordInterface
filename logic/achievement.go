package logic

import (
	"go.uber.org/zap"
	"go_web/dao/mysql"
	"go_web/models"
	"strconv"
	"time"
)

// UpdateAchievement 更新连续天数
func UpdateAchievement(id string) {
	updateTime, _ := mysql.SelUpdateTime(id)
	strTime := updateTime
	t1, _ := time.Parse(time.RFC3339, strTime)
	s1 := t1.Format("2006-01-02 15:04:05")
	daySQL, _ := strconv.Atoi(s1[8:10])
	day := time.Now().Day()
	zap.L().Info("test", zap.Int("daysSQL", daySQL), zap.Int("day", day))
	if day-daySQL == 1 {
		// 如果今天第一次使用，则更新连续天数
		// 获取连续天数
		var achievement models.Achievement
		achievement, _ = mysql.GetAchievement(id)
		achievement.Latest = achievement.Latest + 1
		if achievement.Latest > achievement.Max { // 如果最近连续天数大于最大连续天数，则更新数据库
			achievement.Max = achievement.Latest
		}
		// 将其更新至数据库
		zap.L().Info("更新连续天数", zap.Int("max_achievement", achievement.Max), zap.Int("latest_achievement", achievement.Latest))
		mysql.UpdateAchievement(id, achievement.Max, achievement.Latest)
	} else {
		zap.L().Info("连续天数无需更新", zap.String("id", id))
	}
}

// ClearAchievement 凌晨检测，是否清空最近连续天数
func ClearAchievement() {
	day := time.Now().Day()
	hour := time.Now().Hour()
	minute := time.Now().Minute()
	if hour == 0 && minute == 0 {
		var updateTime []models.UpdateTime
		updateTime = mysql.GetAllUpdateTime()
		for i := range updateTime {
			strTime := updateTime[i].UpdateTime
			t1, _ := time.Parse(time.RFC3339, strTime)
			s1 := t1.Format("2006-01-02 15:04:05")
			daySQL, _ := strconv.Atoi(s1[8:10])
			if day-daySQL > 1 {
				zap.L().Info("清空最近连续天数", zap.String("id", updateTime[i].ID))
				mysql.LatestSetZero(updateTime[i].ID)
			}
		}
		time.Sleep(60 * time.Second)
	}
}
