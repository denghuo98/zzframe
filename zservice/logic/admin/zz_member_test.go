package admin

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zservice"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

// setupTestDBForMember 设置成员测试数据库
func setupTestDBForMember() {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Link: "sqlite::@file(:memory:)?cache=shared",
			},
		},
	})

	// 创建角色表（因为member依赖role）
	createRoleTableSQL := `
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

	// 创建成员表
	createMemberTableSQL := `
	CREATE TABLE zz_admin_member (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		dept_id INTEGER DEFAULT 0,
		real_name VARCHAR(255),
		username VARCHAR(255) NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		salt VARCHAR(255) NOT NULL,
		password_reset_token VARCHAR(255),
		avatar VARCHAR(255),
		sex INTEGER DEFAULT 0,
		email VARCHAR(255),
		mobile VARCHAR(255),
		last_active_at DATETIME,
		remark TEXT,
		status INTEGER DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	ctx := gctx.New()
	_, err := g.DB().Exec(ctx, createRoleTableSQL)
	if err != nil {
		panic(err)
	}

	_, err = g.DB().Exec(ctx, createMemberTableSQL)
	if err != nil {
		panic(err)
	}

	// 创建测试角色（因为member需要关联role）
	testRole := admin.RoleEditInput{
		AdminRole: entity.AdminRole{
			Name:   "测试角色",
			Key:    "test_role",
			Remark: "用于测试的角色",
			Sort:   1,
			Status: 1,
		},
	}
	roleService := &sAdminRole{}
	_, err = roleService.Edit(ctx, testRole)
	if err != nil {
		panic(err)
	}

	// 注册角色服务
	zservice.RegisterAdminRole(roleService)
}

// cleanupTestDBForMember 清理成员测试数据库
func cleanupTestDBForMember() {
	ctx := gctx.New()
	// 删除表
	g.DB().Exec(ctx, "DROP TABLE IF EXISTS zz_admin_member")
	g.DB().Exec(ctx, "DROP TABLE IF EXISTS zz_admin_role")
}

func TestAdminMember_Edit(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 设置测试数据库
		setupTestDBForMember()
		defer cleanupTestDBForMember()

		ctx := gctx.New()
		s := &sAdminMember{}

		// 获取测试角色ID
		testRole, err := dao.AdminRole.Ctx(ctx).Where("name", "测试角色").One()
		t.AssertNil(err)
		roleId := testRole["id"].Int64()

		// 测试新增用户 - 正常情况
		t.Run("CreateMember", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Username: "testuser",
				Password: "123456",
				RealName: "测试用户",
				RoleIds:  []int64{roleId},
				Email:    "test@example.com",
				Mobile:   "13800138000",
				Sex:      1,
				Status:   1,
			}

			err := s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确插入
			record, err := dao.AdminMember.Ctx(ctx).Where("username", "testuser").One()
			t.AssertNil(err)
			t.AssertNE(record, nil)
			t.Assert(record["username"], "testuser")
			t.Assert(record["real_name"], "测试用户")
			t.Assert(record["email"], "test@example.com")
			t.Assert(record["mobile"], "13800138000")
			t.Assert(record["sex"], 1)
			t.Assert(record["status"], 1)
			// 验证密码hash不为空
			t.Assert(record["password_hash"].String() != "", true)
			t.Assert(record["salt"].String() != "", true)
		})

		// 测试新增用户 - 账号为空
		t.Run("CreateMemberEmptyUsername", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Username: "", // 空账号
				Password: "123456",
				RealName: "测试用户2",
				RoleIds:  []int64{roleId},
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "账号不能为空")
		})

		// 测试新增用户 - 用户名重复
		t.Run("CreateMemberDuplicateUsername", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Username: "testuser", // 已存在的用户名
				Password: "123456",
				RealName: "测试用户重复",
				RoleIds:  []int64{roleId},
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "账号已存在，请更换一个")
		})

		// 测试新增用户 - 邮箱重复
		t.Run("CreateMemberDuplicateEmail", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Username: "testuser2",
				Password: "123456",
				RealName: "测试用户邮箱重复",
				RoleIds:  []int64{roleId},
				Email:    "test@example.com", // 已存在的邮箱
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "邮箱已存在，请更换一个")
		})

		// 测试新增用户 - 手机号重复
		t.Run("CreateMemberDuplicateMobile", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Username: "testuser3",
				Password: "123456",
				RealName: "测试用户手机号重复",
				RoleIds:  []int64{roleId},
				Mobile:   "13800138000", // 已存在的手机号
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "手机号码已存在，请更换一个")
		})

		// 测试新增用户 - 角色ID不存在
		t.Run("CreateMemberInvalidRoleId", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Username: "testuser4",
				Password: "123456",
				RealName: "测试用户无效角色",
				RoleIds:  []int64{99999}, // 不存在的角色ID
			}

			err := s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "角色不存在")
		})

		// 获取已创建用户的ID用于更新测试
		existingUser, err := dao.AdminMember.Ctx(ctx).Where("username", "testuser").One()
		t.AssertNil(err)
		userId := existingUser["id"].Int64()

		// 测试更新用户 - 正常更新（不包含密码）
		t.Run("UpdateMember", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Id:       userId,
				Username: "testuser",
				RealName: "测试用户已更新",
				RoleIds:  []int64{roleId},
				Email:    "updated@example.com",
				Mobile:   "13800138001",
				Sex:      2,
				Remark:   "已更新",
			}

			err := s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证数据是否正确更新
			record, err := dao.AdminMember.Ctx(ctx).WherePri(userId).One()
			t.AssertNil(err)
			t.Assert(record["real_name"], "测试用户已更新")
			t.Assert(record["email"], "updated@example.com")
			t.Assert(record["mobile"], "13800138001")
			t.Assert(record["sex"], 2)
			t.Assert(record["remark"], "已更新")
		})

		// 测试更新用户 - 更新密码
		t.Run("UpdateMemberPassword", func(subT *testing.T) {
			input := &admin.MemberEditInput{
				Id:       userId,
				Username: "testuser",
				Password: "newpassword123",
				RealName: "测试用户已更新",
			}

			err := s.Edit(ctx, input)
			t.AssertNil(err)

			// 验证密码是否正确更新
			record, err := dao.AdminMember.Ctx(ctx).WherePri(userId).One()
			t.AssertNil(err)
			// 密码hash应该改变
			t.Assert(record["password_hash"].String() != "", true)
		})

		// 测试更新用户 - 用户名重复
		t.Run("UpdateMemberDuplicateUsername", func(subT *testing.T) {
			// 先创建另一个用户
			input2 := &admin.MemberEditInput{
				Username: "testuser5",
				Password: "123456",
				RealName: "另一个用户",
				RoleIds:  []int64{roleId},
			}
			err := s.Edit(ctx, input2)
			t.AssertNil(err)

			// 尝试将第一个用户更新为与第二个用户相同的用户名
			input := &admin.MemberEditInput{
				Id:       userId,
				Username: "testuser5", // 已存在的用户名
				RealName: "测试用户用户名重复",
			}

			err = s.Edit(ctx, input)
			t.AssertNE(err, nil)
			t.Assert(err.Error(), "账号已存在，请更换一个")
		})
	})
}
