package admin

import (
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zschema/zform"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

// setupTestDB 设置测试数据库
func setupTestDB() {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Link: "sqlite::@file(:memory:)?cache=shared",
			},
		},
	})
	// 创建表结构
	createTableSQL := `
	CREATE TABLE zz_admin_role (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL,
		"key" VARCHAR(255) NOT NULL,
		remark TEXT,
		sort INTEGER DEFAULT 0,
		status INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	ctx := gctx.New()
	_, err := g.DB().Exec(ctx, createTableSQL)
	if err != nil {
		panic(err)
	}

}

// cleanupTestDB 清理测试数据库
func cleanupTestDB() {
	ctx := gctx.New()
	// 删除表
	_, err := g.DB().Exec(ctx, "DROP TABLE IF EXISTS zz_admin_role")
	if err != nil {
		panic(err)
	}
}

func TestAdminRole_RoleEdit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDB()
		defer cleanupTestDB()

		ctx := gctx.New()
		s := &sAdminRole{}

		// 测试新增角色
		t.Run("CreateRole", func(subT *testing.T) {
			input := admin.RoleEditInput{
				AdminRole: entity.AdminRole{
					Name:   "管理员",
					Key:    "admin",
					Remark: "系统管理员",
					Sort:   1,
					Status: 1,
				},
			}

			_, err := s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确插入
			record, err := dao.AdminRole.Ctx(ctx).Where("name", "管理员").One()
			t.AssertNil(err)
			t.AssertNE(record, nil)
			t.Assert(record["name"], "管理员")
			t.Assert(record["key"], "admin")
		})

		// 测试更新角色
		t.Run("UpdateRole", func(subT *testing.T) {
			// 先获取已创建的角色ID
			record, err := dao.AdminRole.Ctx(ctx).Where("name", "管理员").One()
			t.AssertNil(err)
			id := record["id"].Int64()

			input := admin.RoleEditInput{
				AdminRole: entity.AdminRole{
					Id:     id,
					Name:   "超级管理员",
					Key:    "super_admin",
					Remark: "超级管理员角色",
					Sort:   2,
					Status: 1,
				},
			}

			_, err = s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确更新
			updatedRecord, err := dao.AdminRole.Ctx(ctx).WherePri(id).One()
			t.AssertNil(err)
			t.Assert(updatedRecord["name"], "超级管理员")
			t.Assert(updatedRecord["key"], "super_admin")
		})

		// 测试唯一性校验 - 角色名称重复
		t.Run("DuplicateName", func(subT *testing.T) {
			input := admin.RoleEditInput{
				AdminRole: entity.AdminRole{
					Name:   "超级管理员", // 已存在的名称
					Key:    "new_admin",
					Remark: "新管理员",
					Sort:   3,
					Status: 1,
				},
			}

			_, err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "角色名称已存在")
		})

		// 测试唯一性校验 - 角色标识重复
		t.Run("DuplicateKey", func(subT *testing.T) {
			input := admin.RoleEditInput{
				AdminRole: entity.AdminRole{
					Name:   "新管理员",
					Key:    "super_admin", // 已存在的标识
					Remark: "新管理员",
					Sort:   3,
					Status: 1,
				},
			}

			_, err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "角色标识已存在")
		})

		// 测试更新时排除自身记录的唯一性校验
		t.Run("UpdateSelfExclusion", func(subT *testing.T) {
			// 获取现有角色
			record, err := dao.AdminRole.Ctx(ctx).Where("name", "超级管理员").One()
			t.AssertNil(err)
			id := record["id"].Int64()

			// 更新为相同名称和标识（应该成功，因为是更新自身）
			input := admin.RoleEditInput{
				AdminRole: entity.AdminRole{
					Id:     id,
					Name:   "超级管理员",
					Key:    "super_admin",
					Remark: "更新后的超级管理员",
					Sort:   1,
					Status: 1,
				},
			}

			_, err = s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确更新
			updatedRecord, err := dao.AdminRole.Ctx(ctx).WherePri(id).One()
			t.AssertNil(err)
			t.Assert(updatedRecord["remark"], "更新后的超级管理员")
		})
	})
}

func TestAdminRole_Delete(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDB()
		defer cleanupTestDB()

		ctx := gctx.New()
		s := &sAdminRole{}

		// 创建测试数据
		// 创建一个禁用状态的角色用于删除测试
		disabledRole := admin.RoleEditInput{
			AdminRole: entity.AdminRole{
				Name:   "待删除角色",
				Key:    "to_delete",
				Remark: "这个角色将被删除",
				Sort:   10,
				Status: 0, // 禁用状态
			},
		}
		_, err := s.Edit(ctx, disabledRole)
		t.AssertNil(err)

		// 验证角色创建成功且状态正确
		disabledRecordCheck, err := dao.AdminRole.Ctx(ctx).Where("name", "待删除角色").One()
		t.AssertNil(err)
		t.Assert(disabledRecordCheck["status"], 0)

		// 创建一个启用状态的角色用于测试不能删除启用状态的角色
		enabledRole := admin.RoleEditInput{
			AdminRole: entity.AdminRole{
				Name:   "启用角色",
				Key:    "enabled_role",
				Remark: "这个角色是启用状态",
				Sort:   11,
				Status: 1, // 启用状态
			},
		}
		_, err = s.Edit(ctx, enabledRole)
		t.AssertNil(err)

		// 验证启用角色创建成功且状态正确
		enabledRecordCheck, err := dao.AdminRole.Ctx(ctx).Where("name", "启用角色").One()
		t.AssertNil(err)
		t.Assert(enabledRecordCheck["status"], 1)

		// 获取创建的角色ID
		disabledRecord, err := dao.AdminRole.Ctx(ctx).Where("name", "待删除角色").One()
		t.AssertNil(err)
		disabledId := disabledRecord["id"].Int64()

		enabledRecord, err := dao.AdminRole.Ctx(ctx).Where("name", "启用角色").One()
		t.AssertNil(err)
		enabledId := enabledRecord["id"].Int64()

		// 测试正常删除禁用状态的角色
		t.Run("DeleteDisabledRole", func(subT *testing.T) {
			deleteInput := admin.RoleDeleteInput{
				Id: disabledId,
			}

			err := s.Delete(ctx, deleteInput)
			t.AssertNil(err)

			// 验证数据已被删除
			record, err := dao.AdminRole.Ctx(ctx).WherePri(disabledId).One()
			t.AssertNil(err)
			t.Assert(record, nil)
		})

		// 测试删除不存在的角色
		t.Run("DeleteNonExistentRole", func(subT *testing.T) {
			deleteInput := admin.RoleDeleteInput{
				Id: 99999, // 不存在的ID
			}

			err := s.Delete(ctx, deleteInput)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "数据不存在或已经删除")
		})

		// 测试删除启用状态的角色
		t.Run("DeleteEnabledRole", func(subT *testing.T) {
			deleteInput := admin.RoleDeleteInput{
				Id: enabledId,
			}

			err := s.Delete(ctx, deleteInput)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "角色状态为启用，不能删除")

			// 验证数据仍然存在
			record, err := dao.AdminRole.Ctx(ctx).WherePri(enabledId).One()
			t.AssertNil(err)
			t.AssertNE(record, nil)
		})
	})
}

func TestAdminRole_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDB()
		defer cleanupTestDB()

		ctx := gctx.New()
		s := &sAdminRole{}

		// 创建测试数据
		roles := []admin.RoleEditInput{
			{
				AdminRole: entity.AdminRole{
					Name:   "管理员",
					Key:    "admin",
					Remark: "系统管理员",
					Sort:   1,
					Status: 1,
				},
			},
			{
				AdminRole: entity.AdminRole{
					Name:   "编辑员",
					Key:    "editor",
					Remark: "内容编辑员",
					Sort:   2,
					Status: 1,
				},
			},
			{
				AdminRole: entity.AdminRole{
					Name:   "访客",
					Key:    "guest",
					Remark: "普通访客",
					Sort:   3,
					Status: 0, // 禁用状态
				},
			},
			{
				AdminRole: entity.AdminRole{
					Name:   "超级管理员",
					Key:    "super_admin",
					Remark: "超级管理员权限",
					Sort:   4,
					Status: 1,
				},
			},
		}

		// 批量创建测试角色
		for _, role := range roles {
			_, err := s.Edit(ctx, role)
			t.AssertNil(err)
		}

		// 测试基本列表查询（无条件）
		t.Run("ListAll", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 4)
		})

		// 测试按名称模糊查询
		t.Run("ListByName", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
				Name: "管理员",
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 2) // "管理员" 和 "超级管理员"

			// 验证结果包含正确的角色
			names := make([]string, len(result.List))
			for i, role := range result.List {
				names[i] = role.Name
			}
			t.Assert(garray.NewStrArrayFrom(names).Contains("管理员"), true)
			t.Assert(garray.NewStrArrayFrom(names).Contains("超级管理员"), true)
		})

		// 测试按标识模糊查询
		t.Run("ListByKey", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
				Key: "admin",
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 2) // "admin" 和 "super_admin"

			// 验证结果包含正确的角色
			keys := make([]string, len(result.List))
			for i, role := range result.List {
				keys[i] = role.Key
			}
			t.Assert(garray.NewStrArrayFrom(keys).Contains("admin"), true)
			t.Assert(garray.NewStrArrayFrom(keys).Contains("super_admin"), true)
		})

		// 测试按状态精确查询
		t.Run("ListByStatus", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
				Status: 1, // 启用状态
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 3) // 三个启用状态的角色

			// 验证所有结果都是启用状态
			for _, role := range result.List {
				t.Assert(role.Status, 1)
			}
		})

		// 测试组合条件查询
		t.Run("ListByMultipleConditions", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
				Name:   "员", // 模糊匹配包含"员"的角色
				Status: 1,   // 启用状态
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 3) // "管理员"、"编辑员"、"超级管理员"

			// 验证结果
			names := make([]string, len(result.List))
			for i, role := range result.List {
				names[i] = role.Name
				t.Assert(role.Status, 1) // 都是启用状态
			}
			t.Assert(garray.NewStrArrayFrom(names).Contains("管理员"), true)
			t.Assert(garray.NewStrArrayFrom(names).Contains("编辑员"), true)
			t.Assert(garray.NewStrArrayFrom(names).Contains("超级管理员"), true)
		})

		// 测试分页功能
		t.Run("ListWithPagination", func(subT *testing.T) {
			// 第一页，每页2条
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 2,
				},
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 2)

			// 第二页
			input.Page = 2
			result2, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result2.List), 2)

			// 验证两页数据不同
			firstPageIds := make([]int64, len(result.List))
			secondPageIds := make([]int64, len(result2.List))

			for i, role := range result.List {
				firstPageIds[i] = role.Id
			}
			for i, role := range result2.List {
				secondPageIds[i] = role.Id
			}

			// 确保没有重复的ID
			idMap := make(map[int64]bool)
			allIds := append(firstPageIds, secondPageIds...)
			for _, id := range allIds {
				idMap[id] = true
			}
			t.Assert(len(idMap), 4)
		})

		// 测试排序功能（按sort和id升序）
		t.Run("ListWithSorting", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 4)

			// 验证排序：sort升序，id升序
			for i := 0; i < len(result.List)-1; i++ {
				current := result.List[i]
				next := result.List[i+1]

				// 如果sort相同，则id应该升序
				if current.Sort == next.Sort {
					t.Assert(current.Id < next.Id, true)
				} else {
					// 否则sort应该升序
					t.Assert(current.Sort < next.Sort, true)
				}
			}
		})

		// 测试无数据查询
		t.Run("ListWithNoMatch", func(subT *testing.T) {
			input := admin.RoleListInput{
				PageReq: zform.PageReq{
					Page:    1,
					PerPage: 10,
				},
				Name:   "不存在的角色",
				Status: 1,
			}

			result, _, err := s.List(ctx, input)
			t.AssertNil(err)
			t.Assert(len(result.List), 0)
		})
	})
}
