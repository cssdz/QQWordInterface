package mysql

import (
	"errors"
	"go.uber.org/zap"
	"go_web/models"
)

// GetRecordWordsNum 获取复习单词数量
func GetRecordWordsNum(id string) (int, error) {
	var checkId int
	sqlStr := "SELECT COUNT(`no`) FROM `word`.`user_record` WHERE `id` = ?"
	err := db.Get(&checkId, sqlStr, id)
	if err != nil {
		zap.L().Warn("GetReviewWordsNum select err", zap.Error(err), zap.String("id", id))
		return -1, errors.New("获取复习单词数量失败")
	}
	return checkId, nil
}

// GetNextWordNo 获取用户要背的下一个单词编号
func GetNextWordNo(id string) (int, error) {
	var nextWordNo int
	sqlStr := "SELECT `next_word_no` FROM `word`.`user_info` WHERE `id` = ?"
	err := db.Get(&nextWordNo, sqlStr, id)
	if err != nil {
		zap.L().Warn("GetNextWordNo select err", zap.Error(err), zap.String("id", id))
		return -1, errors.New("获取单词失败")
	}
	return nextWordNo, nil
}

// GetNewWords 获取新单词
func GetNewWords(id string, nextNo int, maxNo int) ([]models.Word, error) {
	var words []models.Word
	sqlStr := "SELECT `no`, `word`, `phone`, `translation`, `sentence`, `phrase`, `word_morph` FROM `word`.vocabulary_kaoyanluan WHERE `no` >= ? and `no` <= ?"
	err := db.Select(&words, sqlStr, nextNo, maxNo)
	if err != nil {
		zap.L().Error("GetReviewWordsNum select err", zap.Error(err), zap.String("id", id))
		return words, errors.New("获取单词失败")
	}
	return words, nil
}

// GetRecordWords 获取记录的单词
func GetRecordWords(id string, num int) ([]models.Word, error) {
	var words []models.Word
	sqlStr := "SELECT `no`, `word`, `phone`, `translation`, `sentence`, `phrase`, `word_morph` FROM `word`.vocabulary_kaoyanluan WHERE `no` IN (SELECT t.no FROM (SELECT `no` FROM `word`.user_record WHERE `id` = ? ORDER BY `count` DESC LIMIT ?)as t)"
	err := db.Select(&words, sqlStr, id, num)
	if err != nil {
		zap.L().Error("GetRecordWords select err", zap.Error(err), zap.String("id", id))
		return words, errors.New("获取单词失败")
	}
	return words, nil
}

// GetUserNum 获取用户默认单词数
func GetUserNum(id string) (int, error) {
	var num int
	sqlStr := "SELECT `num` FROM `word`.user_info WHERE `id` = ?"
	err := db.Get(&num, sqlStr, id)
	if err != nil {
		zap.L().Error("GetUserNum get err", zap.Error(err), zap.String("id", id))
		return num, errors.New("获取单词失败")
	}
	return num, nil
}

// UseTimesAddOne 用户请求记录数+1
func UseTimesAddOne(id string) error {
	sqlStr := "UPDATE `word`.`user_info` SET  `use_times` = `use_times` + 1, `update_time` = CURRENT_TIME WHERE `id` = ?"
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		zap.L().Error("UseTimesAddOne exec err", zap.Error(err), zap.String("id", id))
		return errors.New("单词请求失败")
	}
	return nil
}

// SelUpdateTime 查询更新时间
func SelUpdateTime(id string) (string, error) {
	var updateTime string
	sqlStr := "SELECT `update_time` FROM `word`.user_info WHERE `id` = ?"
	err := db.Get(&updateTime, sqlStr, id)
	if err != nil {
		zap.L().Info("SelUpdateTime get err", zap.Error(err), zap.String("id", id))
		return "", errors.New("未查询到更新时间")
	}
	return updateTime, nil
}

// GetAchievement 获取连续天数
func GetAchievement(id string) (models.Achievement, error) {
	var achievement models.Achievement
	sqlStr := "SELECT `max_achievement`, `latest_achievement` FROM `word`.user_info WHERE `id` = ?"
	err := db.Get(&achievement, sqlStr, id)
	if err != nil {
		zap.L().Info("SelUpdateTime get err", zap.Error(err), zap.String("id", id))
		return achievement, errors.New("未查询到更新时间")
	}
	return achievement, nil
}

// UpdateAchievement 更新连续天数
func UpdateAchievement(id string, achievement models.Achievement) {
	sqlStr := "UPDATE `word`.`user_info` SET `max_achievement` = ?, `latest_achievement` = ? WHERE `id` = ?"
	_, err := db.Exec(sqlStr, achievement.Max, achievement.Latest, id)
	if err != nil {
		zap.L().Error("UpdateAchievement exec err", zap.Error(err), zap.String("id", id))
		return
	}
	return
}

// GetAllUpdateTime 获取所有用户的更新时间
func GetAllUpdateTime() []models.UpdateTime {
	var updateTime []models.UpdateTime
	sqlStr := "SELECT `id`, `update_time` FROM `word`.user_info"
	err := db.Select(&updateTime, sqlStr)
	if err != nil {
		zap.L().Error("GetAllUpdateTime select err", zap.Error(err))
		return updateTime
	}
	return updateTime
}

// LatestSetZero 清空最近连续时间
func LatestSetZero(id string) {
	sqlStr := "UPDATE `word`.`user_info` SET `latest_achievement` = 0 WHERE `id` = ?"
	_, err := db.Exec(sqlStr, id)
	if err != nil {
		zap.L().Error("UpdateAchievement exec err", zap.Error(err), zap.String("id", id))
		return
	}
	return
}
