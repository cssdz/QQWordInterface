package mysql

import (
	"errors"
	"go.uber.org/zap"
	"go_web/models"
)

// WordBookToMySQL 将单词书导入MySQL
func WordBookToMySQL(wordInfo models.WordInfo) error {
	sqlStr := "INSERT INTO word.vocabulary_kaoyanluan (no, word, phone, translation, sentence, phrase, word_morph) VALUE (?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(sqlStr, wordInfo.No, wordInfo.Word, wordInfo.Phone, wordInfo.Translation,
		wordInfo.Sentence, wordInfo.Phrase, wordInfo.WordMorph)
	if err != nil {
		zap.L().Error("WordBookToMySQL exec err", zap.Error(err))
		return errors.New("充次错误")
	}
	return nil
}
