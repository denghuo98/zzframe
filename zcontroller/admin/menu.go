package admin

import (
	adminApi "github.com/denghuo98/zzframe/zapi/admin"
	"github.com/denghuo98/zzframe/zservice"
	"github.com/gogf/gf/v2/frame/g"
)

var Menu = cMenu{}

type cMenu struct{}

// Edit 编辑菜单（新增/更新）
func (c *cMenu) Edit(ctx g.Ctx, req *adminApi.MenuEditReq) (res *adminApi.MenuEditRes, err error) {
	if err = zservice.AdminMenu().Edit(ctx, &req.MenuEditInput); err != nil {
		return nil, err
	}
	return &adminApi.MenuEditRes{}, nil
}

// Delete 删除菜单
func (c *cMenu) Delete(ctx g.Ctx, req *adminApi.MenuDeleteReq) (res *adminApi.MenuDeleteRes, err error) {
	if err = zservice.AdminMenu().Delete(ctx, &req.MenuDeleteInput); err != nil {
		return nil, err
	}
	return &adminApi.MenuDeleteRes{}, nil
}

// List 获取菜单列表
func (c *cMenu) List(ctx g.Ctx, req *adminApi.MenuListReq) (res *adminApi.MenuListRes, err error) {
	list, err := zservice.AdminMenu().List(ctx, &req.MenuListInput)
	if err != nil {
		return nil, err
	}
	res = new(adminApi.MenuListRes)
	res.MenuListOutput = *list
	return res, nil
}

// Dynamic 获取动态菜单（用于前端菜单渲染）
func (c *cMenu) Dynamic(ctx g.Ctx, req *adminApi.MenuDynamicReq) (res *adminApi.MenuDynamicRes, err error) {
	list, err := zservice.AdminMenu().GetDynamicMenus(ctx)
	if err != nil {
		return nil, err
	}
	res = new(adminApi.MenuDynamicRes)
	res.List = list
	return res, nil
}
