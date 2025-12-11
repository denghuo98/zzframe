package admin

import (
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zschema/zform"
	"github.com/gogf/gf/v2/frame/g"
)

type RoleEditReq struct {
	g.Meta `path:"/role/edit" method:"post" tags:"SYS-03-角色管理" summary:"修改/新增角色"`
	adminSchema.RoleEditInput
}

type RoleEditRes struct {
}

type RoleDeleteReq struct {
	g.Meta `path:"/role/delete" method:"post" tags:"SYS-03-角色管理" summary:"删除角色"`
	adminSchema.RoleDeleteInput
}

type RoleDeleteRes struct {
}

type RoleListReq struct {
	g.Meta `path:"/role/list" method:"get" tags:"SYS-03-角色管理" summary:"获取角色列表"`
	adminSchema.RoleListInput
}

type RoleListRes struct {
	*adminSchema.RoleListOutput
	zform.PageRes
}
