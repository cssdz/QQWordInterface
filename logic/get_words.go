package logic

import (
	"fmt"
	"go_web/models"
)

// Combine 将记录的单词切片和新单词的切片拼接在一起
func Combine(recordWords []models.Word, NewWords []models.Word) []models.Word {
	for _, v := range NewWords {
		recordWords = append(recordWords, v)
	}
	return recordWords
}

// ToString 将单词信息切片转换成返回请求的格式
func ToString(words []models.Word) []models.WordString {
	var WordStringExport []models.WordString
	for i := range words {
		export := models.WordString{
			No:   words[i].No,
			Word: words[i].Words,
			Content: fmt.Sprintf("%s\t[%s]\n%s\n%s\n%s\n%s", words[i].Words, words[i].Phone,
				words[i].Translation, words[i].Sentence, words[i].Phrase, words[i].WordMorph),
		}
		WordStringExport = append(WordStringExport, export)
	}
	return WordStringExport
}
