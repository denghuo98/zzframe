package admin

import (
	"fmt"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/denghuo98/zzframe/zdb/zgorm"
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
	"github.com/denghuo98/zzframe/zservice"
)

type sMenu struct{}

func init() {
	// 设置单例实例
	zservice.RegisterAdminMenu(&sMenu{})
}

func (s *sMenu) VerifyUnique(ctx g.Ctx, in *adminSchema.VerifyUniqueInput) error {
	if in.Where == nil {
		return nil
	}

	cols := dao.AdminMenu.Columns()
	msgMap := g.MapStrStr{
		cols.Name:  "菜单名称已存在，请换一个",
		cols.Title: "菜单标题已存在，请换一个",
	}

	for k, v := range in.Where {
		if v == "" {
			continue
		}
		msg, ok := msgMap[k]
		if !ok {
			return gerror.Newf("字段 [%v] 未配置唯一属性验证", k)
		}
		if err := zgorm.IsUnique(ctx, dao.AdminMenu, g.Map{k: v}, msg, in.Id); err != nil {
			return err
		}
	}
	return nil
}

func (s *sMenu) Edit(ctx g.Ctx, in *adminSchema.MenuEditInput) error {
	// 验证唯一性
	if err := s.VerifyUnique(ctx, &adminSchema.VerifyUniqueInput{
		Id: in.Id,
		Where: g.Map{
			dao.AdminMenu.Columns().Name:  in.Name,
			dao.AdminMenu.Columns().Title: in.Title,
		},
	}); err != nil {
		return err
	}

	if in.Id > 0 {
		if _, err := dao.AdminMenu.Ctx(ctx).WherePri(in.Id).Data(in).Update(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		}
	} else {
		if _, err := dao.AdminMenu.Ctx(ctx).Data(in).OmitEmpty().Insert(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		}
	}
	return nil
}

func (s *sMenu) Delete(ctx g.Ctx, in *adminSchema.MenuDeleteInput) error {
	if in.Id <= 0 {
		return gerror.New("ID不能为空")
	}
	ok, err := dao.AdminMenu.Ctx(ctx).Where(dao.AdminMenu.Columns().Pid, in.Id).Exist()
	if err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if ok {
		return gerror.New("该菜单下有子菜单，不能删除")
	}
	if _, err := dao.AdminMenu.Ctx(ctx).WherePri(in.Id).Delete(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	return nil
}

func (s *sMenu) List(ctx g.Ctx, in *adminSchema.MenuListInput) (*adminSchema.MenuListOutput, error) {
	var user = zcontext.GetUser(ctx)
	if user == nil {
		return nil, gerror.New("用户信息不存在")
	}

	var models []*entity.AdminMenu

	m := dao.AdminMenu.Ctx(ctx)
	cols := dao.AdminMenu.Columns()

	if in.Name != "" {
		m = m.WhereLike(cols.Name, "%"+in.Name+"%")
	}

	// 如果不是超级管理员，需要根据权限来设置菜单的可见性
	if !zservice.AdminMember().VerifySuperAdmin(ctx, user.Id) {
		roleIds := lo.Map(user.Roles, func(role *webSchema.IdentityRole, _ int) int64 {
			return role.Id
		})
		menuIds, err := dao.AdminRoleMenu.Ctx(ctx).WhereIn(dao.AdminRoleMenu.Columns().RoleId, roleIds).Array(dao.AdminRoleMenu.Columns().MenuId)
		if err != nil {
			return nil, gerror.Wrap(err, zconsts.ErrorORM)
		}
		m = m.WhereIn(cols.Id, menuIds)
	}

	// 排序
	orderBy := fmt.Sprintf("%s asc, %s desc", cols.Sort, cols.Id)
	if err := m.Order(orderBy).Scan(&models); err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}

	// 按照父子结构组成树形结构
	treeList := s.treeList(0, models)
	return &adminSchema.MenuListOutput{List: treeList}, nil
}

func (s *sMenu) treeList(pid int64, menuList []*entity.AdminMenu) []*adminSchema.MenuTree {
	var treeList []*adminSchema.MenuTree
	for _, menu := range menuList {
		if menu.Pid == pid {
			tree := &adminSchema.MenuTree{
				AdminMenu: *menu,
			}
			tree.Children = s.treeList(menu.Id, menuList)
			treeList = append(treeList, tree)
		}
	}
	return treeList
}

// GetDynamicMenus 获取动态菜单 - 用于前端菜单渲染
func (s *sMenu) GetDynamicMenus(ctx g.Ctx) ([]*adminSchema.MenuDynamicItem, error) {
	allMenus, err := s.List(ctx, &adminSchema.MenuListInput{})
	if err != nil {
		return nil, err
	}

	return s.generateNuxtUIMenus(allMenus.List), nil
}

func (s *sMenu) generateNuxtUIMenus(menuList []*adminSchema.MenuTree) []*adminSchema.MenuDynamicItem {
	if len(menuList) == 0 {
		return []*adminSchema.MenuDynamicItem{}
	}

	dynamicMenus := make([]*adminSchema.MenuDynamicItem, 0)
	for _, menu := range menuList {
		dynamicMenus = append(dynamicMenus, &adminSchema.MenuDynamicItem{
			Name:      menu.Name,
			Path:      menu.Path,
			Icon:      menu.Icon,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: adminSchema.MenuDynamicItemMeta{
				Label:  menu.Title,
				Icon:   menu.Icon,
				Hidden: menu.Hidden == 1,
				Sort:   menu.Sort,
			},
			Children: s.generateNuxtUIMenus(menu.Children),
		})
	}
	return dynamicMenus
}
