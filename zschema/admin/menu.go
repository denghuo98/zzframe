package admin

import (
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/zform"
	"github.com/gogf/gf/v2/errors/gerror"
)

// MenuEditInput 服务层-编辑菜单输入参数
type MenuEditInput struct {
	entity.AdminMenu
}

func (in *MenuEditInput) Validate() error {
	if in.Title == "" {
		return gerror.New("菜单名称不能为空")
	}

	if in.Type != 3 && in.Path == "" {
		return gerror.New("路由地址不能为空")
	}

	if in.Name == "" {
		return gerror.New("路由别名不能为空")
	}
	return nil
}

type MenuDeleteInput struct {
	Id int64 `json:"id" v:"required#ID不能为空"`
}

type MenuListInput struct {
	zform.PageReq
	Name string `json:"name"`
}

type MenuListOutput struct {
	List []*MenuTree `json:"list"`
}

type MenuTree struct {
	entity.AdminMenu
	Children []*MenuTree `json:"children"`
}

type MenuDynamicItemMeta struct {
	Label  string `json:"label" dc:"菜单名称"`
	Icon   string `json:"icon,omitempty" dc:"菜单图标"`
	Hidden bool   `json:"hidden" dc:"是否隐藏"`
	Sort   int    `json:"sort,omitempty" dc:"排序,越小越靠前"`
}

type MenuDynamicItem struct {
	Name      string              `json:"name" dc:"菜单名称"`
	Path      string              `json:"path" dc:"路由地址"`
	Icon      string              `json:"icon" dc:"图标"`
	Component string              `json:"component" dc:"组件路径"`
	Redirect  string              `json:"redirect" dc:"重定向地址"`
	Meta      MenuDynamicItemMeta `json:"meta" dc:"菜单元数据"`
	Children  []*MenuDynamicItem  `json:"children" dc:"子菜单"`
}
