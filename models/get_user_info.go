package models

// GetUserInfo 获取用户信息
type GetUserInfo struct {
	Name              string `json:"name" db:"name"`
	Num               int    `json:"num,string" db:"num"`
	UseTimes          string `json:"use_times" db:"use_times"`
	WordTimes         string `json:"word_times" db:"word_times"`
	RecordWordNum     string `json:"record_word_num" db:"record_word_num"`
	MaxAchievement    string `json:"max_achievement" db:"max_achievement"`
	LatestAchievement string `json:"latest_achievement" db:"latest_achievement"`
	UpdateTime        string `json:"update_time" db:"update_time"`
	RegisterTime      string `json:"register_time" db:"register_time"`
}

// UserInfoRequest 用户信息请求结构体
type UserInfoRequest struct {
	ID string `json:"id"`
}
