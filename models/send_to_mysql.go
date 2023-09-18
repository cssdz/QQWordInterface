package models

// WordInfo 定义单词信息结构体
type WordInfo struct {
	No          int    `json:"no" db:"no"`
	Word        string `json:"word" db:"word"`
	Phone       string `json:"phone" db:"phone"`
	Translation string `json:"translation" db:"translation"`
	Sentence    string `json:"sentence" db:"sentence"`
	Phrase      string `json:"phrase" db:"phrase"`
	WordMorph   string `json:"word_morph" db:"word_morph"`
	Comments    string `json:"comments" db:"comments"`
}

// WordName 定义导入单词书名结构体
type WordName struct {
	Name string `json:"name"`
}
