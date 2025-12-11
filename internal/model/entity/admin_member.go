// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// AdminMember is the golang structure for table admin_member.
type AdminMember struct {
	Id                 int64       `json:"id"                 orm:"id"                   description:"管理员ID"`
	DeptId             int64       `json:"deptId"             orm:"dept_id"              description:"部门ID"`
	RealName           string      `json:"realName"           orm:"real_name"            description:"真实姓名"`
	Username           string      `json:"username"           orm:"username"             description:"帐号"`
	PasswordHash       string      `json:"passwordHash"       orm:"password_hash"        description:"密码"`
	Salt               string      `json:"salt"               orm:"salt"                 description:"密码盐"`
	PasswordResetToken string      `json:"passwordResetToken" orm:"password_reset_token" description:"密码重置令牌"`
	Avatar             string      `json:"avatar"             orm:"avatar"               description:"头像"`
	Sex                int         `json:"sex"                orm:"sex"                  description:"性别"`
	Email              string      `json:"email"              orm:"email"                description:"邮箱"`
	Mobile             string      `json:"mobile"             orm:"mobile"               description:"手机号码"`
	LastActiveAt       *gtime.Time `json:"lastActiveAt"       orm:"last_active_at"       description:"最后活跃时间"`
	Remark             string      `json:"remark"             orm:"remark"               description:"备注"`
	Status             int         `json:"status"             orm:"status"               description:"状态"`
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:"创建时间"`
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:"修改时间"`
}
