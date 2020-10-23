package common

// User 用户导出对象
type User struct {
	UserID   int    `json:"userId"`
	UserName string `json:"userName"`
	Status   int    `json:"status"`
}
