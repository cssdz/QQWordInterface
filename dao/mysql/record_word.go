package mysql

import (
	"go.uber.org/zap"
)

// GetNoByWord 获取单词编号
func GetNoByWord(word string) (string, error) {
	var no string
	sqlStr := "SELECT `no` FROM `word`.vocabulary_kaoyanluan WHERE `word` = ?"
	err := db.Get(&no, sqlStr, word)
	if err != nil {
		zap.L().Error("GetAllUpdateTime select err", zap.Error(err))
		return no, err
	}
	return no, nil
}

// GetNoInUserRecord 从用户记录中获取单词编号
func GetNoInUserRecord(no string) error {
	sqlStr := "SELECT `no` FROM `word`.user_record WHERE `no` = ?"
	err := db.Get(&no, sqlStr, no)
	if err != nil {
		zap.L().Info("GetAllUpdateTime select err", zap.Error(err))
		return err
	}
	return nil
}
