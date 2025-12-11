package admin

import "github.com/gogf/gf/v2/frame/g"

// VerifyUniqueInput 验证系统用户唯一性输入
type VerifyUniqueInput struct {
	Id    int64
	Where g.Map
}
