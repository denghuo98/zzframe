package admin

import (
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zschema/zform"
	"github.com/gogf/gf/v2/frame/g"
)

// MemberInfoReq 获取登录用户信息
type MemberInfoReq struct {
	g.Meta `path:"/member/info" method:"get" tags:"SYS-02-用户管理" summary:"获取登录用户信息"`
}

type MemberInfoRes struct {
	*adminSchema.LoginMemberInfoOutput
}

type MemberListReq struct {
	g.Meta `path:"/member/list" method:"get" tags:"SYS-02-用户管理" summary:"获取用户列表"`
	adminSchema.MemberListInput
}

type MemberListRes struct {
	adminSchema.MemberListOutput
	zform.PageRes
}

// EditReq 修改/新增
type MemberEditReq struct {
	g.Meta `path:"/member/edit" method:"post" tags:"SYS-02-用户管理" summary:"修改/新增用户"`
	adminSchema.MemberEditInput
}

type MemberEditRes struct {
}

// DeleteReq 删除用户
type DeleteReq struct {
	g.Meta `path:"/member/delete" method:"post" tags:"SYS-02-用户管理" summary:"删除用户"`
	adminSchema.MemberDeleteInput
}

type DeleteRes struct {
}

// ------------- 以下是个人用户修改自己信息 -------------

// UpdateProfileReq 更新用户资料
type UpdateProfileReq struct {
	g.Meta `path:"/member/update-profile" method:"post" tags:"SYS-02-用户管理" summary:"更新用户资料"`
	adminSchema.MemberUpdateProfileInput
}

type UpdateProfileRes struct {
}

// UpdatePwdReq 更新用户密码
type UpdatePwdReq struct {
	g.Meta `path:"/member/update-pwd" method:"post" tags:"SYS-02-用户管理" summary:"更新用户密码"`
	adminSchema.MemberUpdatePasswordInput
}

type UpdatePwdRes struct {
}
