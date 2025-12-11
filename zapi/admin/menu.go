package admin

import (
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	"github.com/gogf/gf/v2/frame/g"
)

// MenuEditReq 编辑菜单（新增/更新）
type MenuEditReq struct {
	g.Meta `path:"/menu/edit" method:"post" tags:"SYS-01-菜单管理" summary:"编辑菜单（新增/更新）"`
	adminSchema.MenuEditInput
}

type MenuEditRes struct {
}

// MenuDeleteReq 删除菜单
type MenuDeleteReq struct {
	g.Meta `path:"/menu/delete" method:"post" tags:"SYS-01-菜单管理" summary:"删除菜单"`
	adminSchema.MenuDeleteInput
}

type MenuDeleteRes struct {
}

// MenuListReq 获取菜单列表
type MenuListReq struct {
	g.Meta `path:"/menu/list" method:"get" tags:"SYS-01-菜单管理" summary:"获取菜单列表"`
	adminSchema.MenuListInput
}

type MenuListRes struct {
	adminSchema.MenuListOutput
}

// MenuDynamicReq 获取动态菜单（用于前端菜单渲染）
type MenuDynamicReq struct {
	g.Meta `path:"/menu/dynamic" method:"get" tags:"SYS-01-菜单管理" summary:"获取动态菜单"`
}

type MenuDynamicRes struct {
	List []*adminSchema.MenuDynamicItem `json:"list" dc:"动态菜单列表"`
}
