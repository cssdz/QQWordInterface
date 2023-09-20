package mysql

import (
	"errors"
	"go.uber.org/zap"
)

// SelUserRecord 查询用户记录
func SelUserRecord(id string, no string) error {
	var noSql string
	sqlStr := "SELECT `no` FROM `word`.user_record WHERE `id` = ? AND `no` = ?"
	err := db.Get(&noSql, sqlStr, id, no)
	if err != nil {
		zap.L().Info("SelUserRecord get err", zap.Error(err), zap.String("id", id))
		return errors.New("未查询到用户记录")
	}
	return nil
}

// InsertUserRecord 插入用户记录
func InsertUserRecord(id string, no string) error {
	sqlStr := "INSERT INTO `word`.`user_record` (`id`, `no`) VALUE (?, ?)"
	_, err := db.Exec(sqlStr, id, no)
	if err != nil {
		zap.L().Error("InsertUserRecord exec err", zap.Error(err), zap.String("id", id))
		return errors.New("用户记录更新失败")
	}
	return nil
}

// UpdateUserRecord 更新用户记录
func UpdateUserRecord(id string, no string) error {
	sqlStr := "UPDATE `word`.`user_record` SET `count` = `count`+1 WHERE `id` = ? AND `no` = ?"
	_, err := db.Exec(sqlStr, id, no)
	if err != nil {
		zap.L().Error("UpdateUserRecord exec err", zap.Error(err), zap.String("id", id))
		return errors.New("用户记录更新失败")
	}
	return nil
}

// UpdateUser 更新用户信息
func UpdateUser(id string, no string, status int, turn int) error {
	var SqlStr string
	if turn == 0 {
		switch status {
		case 0:
			SqlStr = "UPDATE `word`.`user_info` SET `word_times` = `word_times`+1, `next_word_no` = ?+1, `record_word_num`= `record_word_num`+1, `update_time` = CURRENT_TIME WHERE `id` = ?"
		case 1:
			SqlStr = "UPDATE `word`.`user_info` SET `word_times` = `word_times`+1, `next_word_no` = ?+1, `update_time` = CURRENT_TIME WHERE `id` = ?"
		}
		_, err := db.Exec(SqlStr, no, id)
		if err != nil {
			zap.L().Error("UpdateUserRecord1 exec err", zap.Error(err), zap.String("id", id))
			return errors.New("用户记录更新失败")
		}
	} else {
		switch status {
		case 0:
			SqlStr = "UPDATE `word`.`user_info` SET `word_times` = `word_times`+1, `record_word_num`= `record_word_num`+1, `update_time` = CURRENT_TIME WHERE `id` = ?"
		case 1:
			SqlStr = "UPDATE `word`.`user_info` SET `word_times` = `word_times`+1, `update_time` = CURRENT_TIME WHERE `id` = ?"
		}
		_, err := db.Exec(SqlStr, id)
		if err != nil {
			zap.L().Error("UpdateUserRecord2 exec err", zap.Error(err), zap.String("id", id))
			return errors.New("用户记录更新失败")
		}
	}
	return nil
}
