package models

// Register 注册/设置结构体
type Register struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Num  int    `json:"num,string" db:"num"`
}
