// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminMember is the golang structure of table zz_admin_member for DAO operations like Where/Data.
type AdminMember struct {
	g.Meta             `orm:"table:zz_admin_member, do:true"`
	Id                 any         // 管理员ID
	DeptId             any         // 部门ID
	RealName           any         // 真实姓名
	Username           any         // 帐号
	PasswordHash       any         // 密码
	Salt               any         // 密码盐
	PasswordResetToken any         // 密码重置令牌
	Avatar             any         // 头像
	Sex                any         // 性别
	Email              any         // 邮箱
	Mobile             any         // 手机号码
	LastActiveAt       *gtime.Time // 最后活跃时间
	Remark             any         // 备注
	Status             any         // 状态
	CreatedAt          *gtime.Time // 创建时间
	UpdatedAt          *gtime.Time // 修改时间
}
