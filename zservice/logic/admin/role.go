package admin

import (
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zcasbin"
	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/denghuo98/zzframe/zdb/zgorm"
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zservice"
)

// sAdminRole 后台管理系统中的角色管理
type sAdminRole struct {
}

func init() {
	// 设置单例实例
	zservice.RegisterAdminRole(&sAdminRole{})
}

func (s *sAdminRole) VerifyRoleId(ctx g.Ctx, id int64) (err error) {
	if id <= 0 {
		return gerror.New("角色ID不能为空")
	}

	if ok, err := dao.AdminRole.Ctx(ctx).WherePri(id).Exist(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	} else if !ok {
		return gerror.New("角色不存在")
	}
	return nil
}

func (s *sAdminRole) Verify(ctx g.Ctx, path, method string) bool {
	user := zcontext.GetUser(ctx)
	if user == nil {
		g.Log().Warning(ctx, "用户信息不存在")
		return false
	}

	// 如果用户是超级管理员，则直接返回 true
	if zservice.AdminMember().VerifySuperAdmin(ctx, user.Id) {
		return true
	}

	for _, role := range user.Roles {
		ok, err := zcasbin.Enforcer.Enforce(role.Key, path, method)
		if err != nil {
			g.Log().Warning(ctx, "Casbin Enforcer 验证权限失败", err)
			continue
		}
		if ok {
			return true
		}
	}
	return false

}

func (s *sAdminRole) Edit(ctx g.Ctx, in adminSchema.RoleEditInput) (role *entity.AdminRole, err error) {
	cols := dao.AdminRole.Columns()
	// 判断角色名称是否唯一
	if err = zgorm.IsUnique(ctx, &dao.AdminRole, g.Map{cols.Name: in.Name}, "角色名称已存在", in.Id); err != nil {
		return nil, err
	}

	// 判断角色标识是否唯一
	if err = zgorm.IsUnique(ctx, &dao.AdminRole, g.Map{cols.Key: in.Key}, "角色标识已存在", in.Id); err != nil {
		return nil, err
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	m := g.Model(dao.AdminRole.Table()).TX(tx)
	if in.Id > 0 {
		// 更新角色
		_, err = m.Fields(adminSchema.RoleUpdateFields{}).WherePri(in.Id).Data(in).Update()
		if err != nil {
			return nil, gerror.Wrap(err, zconsts.ErrorORM)
		}
	} else {
		// 新增角色
		in.Id, err = m.Fields(adminSchema.RoleUpdateFields{}).Data(in).OmitEmptyData().InsertAndGetId()
		if err != nil {
			return nil, gerror.Wrap(err, zconsts.ErrorORM)
		}
	}

	// 更新角色绑定的菜单
	if len(in.MenuIds) > 0 {
		if err = s.updateMenus(tx, &adminSchema.RoleUpdateMenuInput{
			Id:      in.Id,
			MenuIds: in.MenuIds,
		}); err != nil {
			return nil, err
		}
	}

	// 重新查询角色
	if err = g.Model(dao.AdminRole.Table()).TX(tx).WherePri(in.Id).Scan(&role); err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}
	return role, nil
}

func (s *sAdminRole) Delete(ctx g.Ctx, in adminSchema.RoleDeleteInput) (err error) {
	var models *entity.AdminRole
	if err = dao.AdminRole.Ctx(ctx).WherePri(in.Id).Scan(&models); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if models == nil {
		return gerror.New("数据不存在或已经删除")
	}
	if models.Status == 1 {
		return gerror.New("角色状态为启用，不能删除")
	}

	if _, err = dao.AdminRole.Ctx(ctx).WherePri(in.Id).Delete(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	return nil
}

func (s *sAdminRole) List(ctx g.Ctx, in adminSchema.RoleListInput) (out *adminSchema.RoleListOutput, totalCount int, err error) {
	var (
		m      = dao.AdminRole.Ctx(ctx)
		models []*entity.AdminRole
		cols   = dao.AdminRole.Columns()
	)

	if in.Name != "" {
		m = m.WhereLike(cols.Name, "%"+in.Name+"%")
	}
	if in.Key != "" {
		m = m.WhereLike(cols.Key, "%"+in.Key+"%")
	}
	if in.Status != 0 {
		m = m.Where(cols.Status, in.Status)
	}

	totalCount, err = m.Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, zconsts.ErrorORM)
	}

	if err = m.Page(in.Page, in.PerPage).Order(fmt.Sprintf("%s asc, %s asc", cols.Sort, cols.Id)).Scan(&models); err != nil {
		return nil, 0, gerror.Wrap(err, zconsts.ErrorORM)
	}
	out = new(adminSchema.RoleListOutput)
	out.List = models
	return out, totalCount, nil
}

// updateMenus 更新角色绑定的菜单
// 参数:
// - tx: 数据库事务
// - in: 更新菜单入参
// 返回:
// - error: 错误信息
func (s *sAdminRole) updateMenus(tx gdb.TX, in *adminSchema.RoleUpdateMenuInput) (err error) {
	// 校验角色是否存在
	if ok, err := g.Model(dao.AdminRole.Table()).TX(tx).Where(dao.AdminRole.Columns().Id, in.Id).Exist(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	} else if !ok {
		return gerror.New("角色不存在")
	}

	// 查询出来当前角色所绑定的菜单ID
	rmCols := dao.AdminRoleMenu.Columns()
	menuIds, err := g.Model(dao.AdminRoleMenu.Table()).TX(tx).Where(rmCols.RoleId, in.Id).Array(rmCols.MenuId)
	if err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	existingMenuIds := make([]int64, 0, len(menuIds))
	for _, menuId := range menuIds {
		existingMenuIds = append(existingMenuIds, menuId.Int64())
	}

	// 将当前角色所绑定的菜单ID与新菜单ID进行对比
	needDeleteMenuIds, needInsertMenuIds := lo.Difference(existingMenuIds, in.MenuIds)

	// 删除不再绑定的菜单
	if len(needDeleteMenuIds) > 0 {
		if _, err = g.Model(dao.AdminRoleMenu.Table()).TX(tx).Where(rmCols.RoleId, in.Id).WhereIn(rmCols.MenuId, needDeleteMenuIds).Delete(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		}
	}

	// 新增需要绑定的菜单
	if len(needInsertMenuIds) > 0 {
		data := make([]*entity.AdminRoleMenu, 0, len(needInsertMenuIds))
		for _, menuId := range needInsertMenuIds {
			data = append(data, &entity.AdminRoleMenu{
				RoleId: in.Id,
				MenuId: menuId,
			})
		}
		if _, err = g.Model(dao.AdminRoleMenu.Table()).TX(tx).Data(data).Insert(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		}
	}
	return nil
}
