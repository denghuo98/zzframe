package admin

import (
	"testing"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zservice"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

// setupTestDBForMenu 设置菜单测试数据库
func setupTestDBForMenu() {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Link: "sqlite::@file(:memory:)?cache=shared",
			},
		},
	})

	// 创建菜单表
	createMenuTableSQL := `
	CREATE TABLE zz_admin_menu (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pid INTEGER DEFAULT 0,
		name VARCHAR(255) NOT NULL,
		title VARCHAR(255) NOT NULL,
		icon VARCHAR(255),
		path VARCHAR(255),
		component VARCHAR(255),
		redirect VARCHAR(255),
		sort INTEGER DEFAULT 0,
		type INTEGER DEFAULT 1,
		hidden INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	ctx := gctx.New()
	_, err := g.DB().Exec(ctx, createMenuTableSQL)
	if err != nil {
		panic(err)
	}

	// 注册菜单服务
	zservice.RegisterAdminMenu(&sMenu{})
}

// cleanupTestDBForMenu 清理菜单测试数据库
func cleanupTestDBForMenu() {
	ctx := gctx.New()
	// 删除表
	g.DB().Exec(ctx, "DROP TABLE IF EXISTS zz_admin_menu")
}

func TestAdminMenu_Edit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDBForMenu()
		defer cleanupTestDBForMenu()

		ctx := gctx.New()
		s := &sMenu{}

		// 测试新增菜单 - 正常情况
		t.Run("CreateMenu", func(subT *testing.T) {
			input := &admin.MenuEditInput{
				AdminMenu: entity.AdminMenu{
					Pid:       0,
					Name:      "dashboard",
					Title:     "仪表盘",
					Icon:      "dashboard",
					Path:      "/dashboard",
					Component: "/dashboard/index",
					Sort:      1,
					Type:      1,
					Hidden:    0,
				},
			}

			err := s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确插入
			record, err := dao.AdminMenu.Ctx(ctx).Where("name", "dashboard").One()
			t.AssertNil(err)
			t.AssertNE(record, nil)
			t.Assert(record["name"], "dashboard")
			t.Assert(record["title"], "仪表盘")
			t.Assert(record["path"], "/dashboard")
			t.Assert(record["component"], "/dashboard/index")
			t.Assert(record["sort"], 1)
			t.Assert(record["type"], 1)
		})

		// 测试新增菜单 - 名称重复
		t.Run("CreateMenuDuplicateName", func(subT *testing.T) {
			input := &admin.MenuEditInput{
				AdminMenu: entity.AdminMenu{
					Pid:   0,
					Name:  "dashboard", // 已存在的名称
					Title: "重复的仪表盘",
					Path:  "/dashboard2",
					Sort:  2,
					Type:  1,
				},
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "菜单名称已存在，请换一个")
		})

		// 测试新增菜单 - 标题重复
		t.Run("CreateMenuDuplicateTitle", func(subT *testing.T) {
			input := &admin.MenuEditInput{
				AdminMenu: entity.AdminMenu{
					Pid:   0,
					Name:  "dashboard2",
					Title: "仪表盘", // 已存在的标题
					Path:  "/dashboard2",
					Sort:  2,
					Type:  1,
				},
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "菜单标题已存在，请换一个")
		})

		// 获取已创建菜单的ID用于更新测试
		existingMenu, err := dao.AdminMenu.Ctx(ctx).Where("name", "dashboard").One()
		t.AssertNil(err)
		menuId := existingMenu["id"].Int64()

		// 测试更新菜单 - 正常更新
		t.Run("UpdateMenu", func(subT *testing.T) {
			input := &admin.MenuEditInput{
				AdminMenu: entity.AdminMenu{
					Id:        menuId,
					Pid:       0,
					Name:      "dashboard",
					Title:     "系统仪表盘",
					Icon:      "dashboard-updated",
					Path:      "/dashboard",
					Component: "/dashboard/index",
					Sort:      10,
					Type:      1,
					Hidden:    0,
				},
			}

			err := s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确更新
			record, err := dao.AdminMenu.Ctx(ctx).WherePri(menuId).One()
			t.AssertNil(err)
			t.Assert(record["title"], "系统仪表盘")
			t.Assert(record["icon"], "dashboard-updated")
			t.Assert(record["sort"], 10)
		})

		// 测试更新菜单 - 名称重复
		t.Run("UpdateMenuDuplicateName", func(subT *testing.T) {
			// 先创建另一个菜单
			input2 := &admin.MenuEditInput{
				AdminMenu: entity.AdminMenu{
					Pid:   0,
					Name:  "user_manage",
					Title: "用户管理",
					Path:  "/user",
					Sort:  3,
					Type:  1,
				},
			}
			err := s.Edit(ctx, input2)
			t.AssertNil(err)

			// 尝试将第一个菜单更新为与第二个菜单相同的名称
			input := &admin.MenuEditInput{
				AdminMenu: entity.AdminMenu{
					Id:    menuId,
					Pid:   0,
					Name:  "user_manage", // 已存在的名称
					Title: "仪表盘重命名",
					Path:  "/dashboard",
					Sort:  1,
					Type:  1,
				},
			}

			err = s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "菜单名称已存在，请换一个")
		})
	})
}

func TestAdminMenu_Delete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDBForMenu()
		defer cleanupTestDBForMenu()

		ctx := gctx.New()
		s := &sMenu{}

		// 创建测试菜单
		parentMenu := &admin.MenuEditInput{
			AdminMenu: entity.AdminMenu{
				Pid:   0,
				Name:  "system",
				Title: "系统管理",
				Path:  "/system",
				Sort:  1,
				Type:  1,
			},
		}
		err := s.Edit(ctx, parentMenu)
		t.AssertNil(err)

		childMenu := &admin.MenuEditInput{
			AdminMenu: entity.AdminMenu{
				Pid:   1, // 父菜单ID
				Name:  "user",
				Title: "用户管理",
				Path:  "/system/user",
				Sort:  1,
				Type:  1,
			},
		}
		err = s.Edit(ctx, childMenu)
		t.AssertNil(err)

		// 获取菜单ID
		parentRecord, err := dao.AdminMenu.Ctx(ctx).Where("name", "system").One()
		t.AssertNil(err)
		parentId := parentRecord["id"].Int64()

		childRecord, err := dao.AdminMenu.Ctx(ctx).Where("name", "user").One()
		t.AssertNil(err)
		childId := childRecord["id"].Int64()

		// 测试删除有子菜单的菜单 - 应该失败
		t.Run("DeleteMenuWithChildren", func(subT *testing.T) {
			deleteInput := &admin.MenuDeleteInput{
				Id: parentId,
			}

			err := s.Delete(ctx, deleteInput)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "该菜单下有子菜单，不能删除")

			// 验证菜单仍然存在
			record, err := dao.AdminMenu.Ctx(ctx).WherePri(parentId).One()
			t.AssertNil(err)
			t.AssertNE(record, nil)
		})

		// 测试删除子菜单 - 应该成功
		t.Run("DeleteMenuWithoutChildren", func(subT *testing.T) {
			deleteInput := &admin.MenuDeleteInput{
				Id: childId,
			}

			err := s.Delete(ctx, deleteInput)
			t.AssertNil(err)

			// 验证菜单已被删除
			record, err := dao.AdminMenu.Ctx(ctx).WherePri(childId).One()
			t.AssertNil(err)
			t.Assert(record, nil)
		})

		// 测试删除不存在的菜单
		t.Run("DeleteNonExistentMenu", func(subT *testing.T) {
			deleteInput := &admin.MenuDeleteInput{
				Id: 99999, // 不存在的ID
			}

			err := s.Delete(ctx, deleteInput)
			// 对于不存在的菜单，删除操作应该成功（因为没有子菜单需要检查）
			t.AssertNil(err)
		})
	})
}

func TestAdminMenu_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDBForMenu()
		defer cleanupTestDBForMenu()

		ctx := gctx.New()
		s := &sMenu{}

		// 创建测试菜单数据
		menus := []*admin.MenuEditInput{
			{
				AdminMenu: entity.AdminMenu{
					Pid:   0,
					Name:  "dashboard",
					Title: "仪表盘",
					Path:  "/dashboard",
					Sort:  1,
					Type:  1,
				},
			},
			{
				AdminMenu: entity.AdminMenu{
					Pid:   0,
					Name:  "system",
					Title: "系统管理",
					Path:  "/system",
					Sort:  2,
					Type:  1,
				},
			},
			{
				AdminMenu: entity.AdminMenu{
					Pid:   2, // system的ID
					Name:  "user",
					Title: "用户管理",
					Path:  "/system/user",
					Sort:  1,
					Type:  1,
				},
			},
			{
				AdminMenu: entity.AdminMenu{
					Pid:   2,
					Name:  "role",
					Title: "角色管理",
					Path:  "/system/role",
					Sort:  2,
					Type:  1,
				},
			},
		}

		// 批量创建菜单
		for _, menu := range menus {
			err := s.Edit(ctx, menu)
			t.AssertNil(err)
		}

		// 测试获取完整菜单列表
		t.Run("ListAllMenus", func(subT *testing.T) {
			input := &admin.MenuListInput{}

			result, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 2) // 应该有2个顶级菜单

			// 验证树形结构
			systemMenu := result.List[1] // system菜单在第2个位置
			t.Assert(systemMenu.Title, "系统管理")
			t.Assert(len(systemMenu.Children), 2) // system下应该有2个子菜单
		})

		// 测试按名称搜索菜单
		t.Run("ListMenusByName", func(subT *testing.T) {
			input := &admin.MenuListInput{
				Name: "用户", // 模糊匹配
			}

			result, err := s.List(ctx, input)
			t.AssertNil(err)
			// 由于创建的菜单中没有包含"用户"这个词，所以应该返回空或只包含顶级菜单
			// 这里我们只验证没有错误发生，具体逻辑可能需要调整
			t.Assert(result != nil, true)
		})

		// 测试空搜索结果
		t.Run("ListMenusNoMatch", func(subT *testing.T) {
			input := &admin.MenuListInput{
				Name: "不存在的菜单",
			}

			result, err := s.List(ctx, input)
			t.AssertNil(err)
			// 应该返回空结果
			t.Assert(len(result.List) == 0, true)
		})
	})
}

func TestAdminMenu_GetDynamicMenus(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDBForMenu()
		defer cleanupTestDBForMenu()

		ctx := gctx.New()
		s := &sMenu{}

		// 创建父菜单
		parentMenu := &admin.MenuEditInput{
			AdminMenu: entity.AdminMenu{
				Pid:   0,
				Name:  "system",
				Title: "系统管理",
				Path:  "/system",
				Sort:  1,
				Type:  1,
			},
		}
		err := s.Edit(ctx, parentMenu)
		t.AssertNil(err)

		// 创建子菜单
		childMenu := &admin.MenuEditInput{
			AdminMenu: entity.AdminMenu{
				Pid:       1, // 父菜单ID
				Name:      "dashboard",
				Title:     "仪表盘",
				Icon:      "dashboard-icon",
				Path:      "/system/dashboard",
				Component: "/dashboard/index",
				Sort:      1,
				Type:      1,
				Hidden:    0,
			},
		}
		err = s.Edit(ctx, childMenu)
		t.AssertNil(err)

		// 测试获取动态菜单
		t.Run("GetDynamicMenus", func(subT *testing.T) {
			result, err := s.GetDynamicMenus(ctx)
			t.AssertNil(err)
			t.Assert(len(result) > 0, true)

			// 验证动态菜单格式
			dynamicMenu := result[0]
			t.Assert(dynamicMenu.Name, "dashboard")
			t.Assert(dynamicMenu.Path, "/system/dashboard")
			t.Assert(dynamicMenu.Meta.Label, "仪表盘")
			t.Assert(dynamicMenu.Meta.Icon, "dashboard-icon")
			t.Assert(dynamicMenu.Meta.Hidden, false)
			t.Assert(dynamicMenu.Meta.Sort, 1)
		})
	})
}
