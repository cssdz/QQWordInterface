package models

// GetWordsRequest 获取单词请求结构体
type GetWordsRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Num  int    `json:"num,string"`
}

// Word 单词信息结构体
type Word struct {
	No          string `json:"no" db:"no"`
	Words       string `json:"words" db:"word"`
	Phone       string `json:"phone" db:"phone"`
	Translation string `json:"translation" db:"translation"`
	Sentence    string `json:"sentence" db:"sentence"`
	Phrase      string `json:"phrase" db:"phrase"`
	WordMorph   string `json:"word_morph" db:"word_morph"`
}

// WordString 返回请求的单词格式
type WordString struct {
	No      string `json:"no"`
	Word    string `json:"word"`
	Content string `json:"content"`
}

// UpdateTime 更新时间结构体
type UpdateTime struct {
	ID         string `json:"id" db:"id"`
	UpdateTime string `json:"update_time" db:"update_time"`
}
