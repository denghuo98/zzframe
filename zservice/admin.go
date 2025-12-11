package zservice

import (
	"github.com/gogf/gf/v2/frame/g"

	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/zschema/admin"
)

// IAdminMenu 后台菜单管理接口
// 定义菜单管理服务的所有方法，实现业务逻辑的解耦
type IAdminMenu interface {
	// VerifyUnique 验证菜单唯一性
	VerifyUnique(ctx g.Ctx, in *admin.VerifyUniqueInput) error

	// Edit 编辑菜单（新增或更新）
	// 如果 in.Id > 0 则更新，否则新增
	Edit(ctx g.Ctx, in *admin.MenuEditInput) error

	// Delete 删除菜单
	// 删除前会检查是否有子菜单
	Delete(ctx g.Ctx, in *admin.MenuDeleteInput) error

	// List 查询菜单列表
	// 返回树形结构的菜单数据
	List(ctx g.Ctx, in *admin.MenuListInput) (*admin.MenuListOutput, error)

	// GetDynamicMenus 获取动态菜单
	// 用于前端菜单渲染，生成NuxtUI菜单格式
	GetDynamicMenus(ctx g.Ctx) ([]*admin.MenuDynamicItem, error)
}

// IAdminRole 后台角色管理接口
// 定义角色管理服务的所有方法，实现业务逻辑的解耦
type IAdminRole interface {

	// VerifyRoleId 验证角色ID
	VerifyRoleId(ctx g.Ctx, id int64) error

	// Verify 验证权限
	Verify(ctx g.Ctx, path, method string) bool

	// Edit 编辑角色（新增或更新）
	// 如果 in.Id > 0 则更新，否则新增
	Edit(ctx g.Ctx, in admin.RoleEditInput) (*entity.AdminRole, error)

	// Delete 删除角色
	// 只能删除状态为禁用的角色
	Delete(ctx g.Ctx, in admin.RoleDeleteInput) error

	// List 查询角色列表
	// 支持按名称、标识、状态过滤，支持分页和排序
	List(ctx g.Ctx, in admin.RoleListInput) (*admin.RoleListOutput, int, error)
}

// IAdminMember 后台成员管理接口
// 定义成员管理服务的所有方法，实现业务逻辑的解耦
type IAdminMember interface {
	// VerifyUnique 验证成员唯一性
	// 检查账号、邮箱、手机号是否已存在
	VerifyUnique(ctx g.Ctx, in *admin.VerifyUniqueInput) error

	// VerifySuperAdmin 验证是否是超级管理员
	VerifySuperAdmin(ctx g.Ctx, id int64) bool

	// Info 获取当前登录用户信息
	Info(ctx g.Ctx) (*admin.LoginMemberInfoOutput, error)

	// List 获取用户列表
	// 支持按账号、姓名、邮箱、手机号、状态过滤，支持分页
	List(ctx g.Ctx, in *admin.MemberListInput) (*admin.MemberListOutput, int, error)

	// Edit 编辑成员（新增或更新）
	// 如果 in.Id > 0 则更新，否则新增
	Edit(ctx g.Ctx, in *admin.MemberEditInput) error

	// UpdateProfile 更新用户资料
	// 更新用户的头像、真实姓名、性别等信息
	UpdateProfile(ctx g.Ctx, in *admin.MemberUpdateProfileInput) error

	// UpdatePassword 更新用户密码
	// 验证旧密码后更新为新密码
	UpdatePassword(ctx g.Ctx, in *admin.MemberUpdatePasswordInput) error

	// Delete 删除成员
	// 软删除用户，将用户状态设置为禁用
	Delete(ctx g.Ctx, in *admin.MemberDeleteInput) error
}

// AdminRole 角色管理单例实例
// 通过接口提供服务，实现依赖倒置和单例访问
var (
	localAdminRole   IAdminRole
	localAdminMenu   IAdminMenu
	localAdminMember IAdminMember
)

func AdminRole() IAdminRole {
	if localAdminRole == nil {
		panic("AdminRole is not initialized, please register it first")
	}
	return localAdminRole
}

func RegisterAdminRole(i IAdminRole) {
	localAdminRole = i
}

// AdminMenu 菜单管理单例实例
// 通过接口提供服务，实现依赖倒置和单例访问
func AdminMenu() IAdminMenu {
	if localAdminMenu == nil {
		panic("AdminMenu is not initialized, please register it first")
	}
	return localAdminMenu
}

func RegisterAdminMenu(i IAdminMenu) {
	localAdminMenu = i
}

// AdminMember 成员管理单例实例
// 通过接口提供服务，实现依赖倒置和单例访问
func AdminMember() IAdminMember {
	if localAdminMember == nil {
		panic("AdminMember is not initialized, please register it first")
	}
	return localAdminMember
}

func RegisterAdminMember(i IAdminMember) {
	localAdminMember = i
}
