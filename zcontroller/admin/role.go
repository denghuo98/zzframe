package admin

import (
	"github.com/gogf/gf/v2/frame/g"

	adminApi "github.com/denghuo98/zzframe/zapi/admin"
	"github.com/denghuo98/zzframe/zservice"
)

var Role = cRole{}

type cRole struct{}

// Edit 修改/新增角色
func (c *cRole) Edit(ctx g.Ctx, req *adminApi.RoleEditReq) (res *adminApi.RoleEditRes, err error) {
	_, err = zservice.AdminRole().Edit(ctx, req.RoleEditInput)
	if err != nil {
		return nil, err
	}
	return &adminApi.RoleEditRes{}, nil
}

// Delete 删除角色
func (c *cRole) Delete(ctx g.Ctx, req *adminApi.RoleDeleteReq) (res *adminApi.RoleDeleteRes, err error) {
	if err = zservice.AdminRole().Delete(ctx, req.RoleDeleteInput); err != nil {
		return nil, err
	}
	return &adminApi.RoleDeleteRes{}, nil
}

// List 获取角色列表
func (c *cRole) List(ctx g.Ctx, req *adminApi.RoleListReq) (res *adminApi.RoleListRes, err error) {
	list, totalCount, err := zservice.AdminRole().List(ctx, req.RoleListInput)
	if err != nil {
		return nil, err
	}
	res = new(adminApi.RoleListRes)
	res.RoleListOutput = list
	res.PageRes.Pack(req, int(totalCount))
	return res, nil
}
