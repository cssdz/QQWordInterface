package mysql

import (
	"errors"
	"go.uber.org/zap"
	"go_web/models"
)

// GetUserInfo 获取用户信息
func GetUserInfo(id string) (models.GetUserInfo, error) {
	var UserInfo models.GetUserInfo
	sqlStr := "SELECT `name`, `num`, `use_times`, `word_times`, `record_word_num`, `max_achievement`, `latest_achievement`, `update_time`, `register_time` FROM `word`.user_info WHERE `id` = ?"
	err := db.Get(&UserInfo, sqlStr, id)
	if err != nil {
		zap.L().Warn("GetUserInfo exec err", zap.Error(err), zap.String("id", id))
		return UserInfo, errors.New("获取用户信息失败")
	}
	return UserInfo, nil
}
