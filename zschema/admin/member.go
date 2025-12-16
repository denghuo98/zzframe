package admin

import (
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/zform"
)

// MemberAddInput 服务层-系统用户新增入参
type MemberAddInput struct {
	*MemberEditInput
	Salt string `json:"salt" dc:"密码盐"`
}

// MemberDeleteInput 服务层-系统用户删除入参
type MemberDeleteInput struct {
	Id int64
}

// MemberEditInput 服务层-系统用户编辑入参
type MemberEditInput struct {
	Id           int64
	Username     string  `json:"username"   v:"required#账号不能为空"           dc:"帐号"`
	PasswordHash string  `json:"passwordHash"                                  dc:"密码hash"`
	Password     string  `json:"password"                                      dc:"密码"`
	RealName     string  `json:"realName"                                      dc:"真实姓名"`
	RoleIds      []int64 `json:"roleIds"                                       dc:"角色ID"`
	Avatar       string  `json:"avatar"                                        dc:"头像"`
	Sex          int     `json:"sex"                                           dc:"性别"`
	Email        string  `json:"email"                                         dc:"邮箱"`
	Mobile       string  `json:"mobile"                                        dc:"手机号码"`
	Remark       string  `json:"remark"                                        dc:"备注"`
	Status       int     `json:"status"                                        dc:"状态"`
}

// MemberUpdateProfileInput 服务层-系统用户更新个人信息入参
type MemberUpdateProfileInput struct {
	Id       int64
	Avatar   string `json:"avatar"   v:"required#头像不能为空"     dc:"头像"`
	RealName string `json:"realName"  v:"required#真实姓名不能为空"       dc:"真实姓名"`
	Sex      int    `json:"sex"         dc:"性别"`
}

// MemberUpdatePasswordInput 服务层-系统用户更新密码入参
type MemberUpdatePasswordInput struct {
	Id          int64
	OldPassword string `json:"oldPassword" v:"required#旧密码不能为空" dc:"旧密码"`
	NewPassword string `json:"newPassword" v:"required#新密码不能为空" dc:"新密码"`
}

// MemberResetPwdInput 服务层-系统用户重置密码入参
type MemberResetPwdInput struct {
	Id       int64
	Password string `json:"password" v:"required#密码不能为空" dc:"密码"`
}

// MemberUpdateRoleInput 服务层-系统用户更新角色入参
type MemberUpdateRoleInput struct {
	Id      int64
	RoleIds []int64 `json:"roleIds" dc:"角色ID"`
}

// LoginMemberInfoOutput 登录用户信息输出
type LoginMemberInfoOutput struct {
	Id       int64
	Username string `json:"username" dc:"账号"`
	RealName string `json:"realName" dc:"真实姓名"`
	Avatar   string `json:"avatar" dc:"头像"`
	Sex      int    `json:"sex" dc:"性别"`
	Email    string `json:"email" dc:"邮箱"`
	Mobile   string `json:"mobile" dc:"手机号码"`
}

// MemberListInput 服务层-系统用户列表入参
type MemberListInput struct {
	zform.PageReq
	Username string `json:"username" dc:"账号"`
	RealName string `json:"realName" dc:"真实姓名"`
	Email    string `json:"email" dc:"邮箱"`
	Mobile   string `json:"mobile" dc:"手机号码"`
	Status   int    `json:"status" dc:"状态"`
}

type MemberListOutputItem struct {
	*entity.AdminMember
	RoleIds []int64             `json:"roleIds" dc:"绑定的角色 ID"`
	Roles   []*entity.AdminRole `json:"roles" dc:"绑定的角色列表"`
}

// MemberListOutput 服务层-系统用户列表输出
type MemberListOutput struct {
	List []*MemberListOutputItem `json:"list" dc:"用户列表"`
}
