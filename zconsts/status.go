package zconsts

// 状态码
const (
	StatusALL     int = -1 // 全部状态
	StatusEnabled int = 1  // 启用
	StatusDisable int = 2  // 禁用
	StatusDelete  int = 3  // 已删除
)

var StatusSlice = []int{StatusALL, StatusEnabled, StatusDisable, StatusDelete}
