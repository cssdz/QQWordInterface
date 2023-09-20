package models

// UpdateUserRequest 更新用户记录请求
type UpdateUserRequest struct {
	ID     string `json:"id"`
	No     string `json:"no"`
	Status string `json:"status"`
}

// Achievement 获取连续天数
type Achievement struct {
	Max    int `json:"max" db:"max_achievement"`
	Latest int `json:"latest" db:"latest_achievement"`
}
