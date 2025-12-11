// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AdminMenuDao is the data access object for the table zz_admin_menu.
type AdminMenuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  AdminMenuColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// AdminMenuColumns defines and stores column names for the table zz_admin_menu.
type AdminMenuColumns struct {
	Id          string // 菜单ID
	Pid         string // 父菜单ID
	Title       string // 菜单名称
	Name        string // 名称编码
	Path        string // 路由地址
	Icon        string // 菜单图标
	Type        string // 菜单类型（1目录 2菜单 3按钮）
	Redirect    string // 重定向地址
	Component   string // 组件路径
	IsFrame     string // 是否内嵌
	FrameSrc    string // 内联外部地址
	Hidden      string // 是否隐藏
	Sort        string // 排序
	Remark      string // 备注
	Status      string // 菜单状态
	UpdatedAt   string // 更新时间
	CreatedAt   string // 创建时间
	Permissions string // 菜单权限
}

// adminMenuColumns holds the columns for the table zz_admin_menu.
var adminMenuColumns = AdminMenuColumns{
	Id:          "id",
	Pid:         "pid",
	Title:       "title",
	Name:        "name",
	Path:        "path",
	Icon:        "icon",
	Type:        "type",
	Redirect:    "redirect",
	Component:   "component",
	IsFrame:     "is_frame",
	FrameSrc:    "frame_src",
	Hidden:      "hidden",
	Sort:        "sort",
	Remark:      "remark",
	Status:      "status",
	UpdatedAt:   "updated_at",
	CreatedAt:   "created_at",
	Permissions: "permissions",
}

// NewAdminMenuDao creates and returns a new DAO object for table data access.
func NewAdminMenuDao(handlers ...gdb.ModelHandler) *AdminMenuDao {
	return &AdminMenuDao{
		group:    "default",
		table:    "zz_admin_menu",
		columns:  adminMenuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AdminMenuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AdminMenuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AdminMenuDao) Columns() AdminMenuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AdminMenuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AdminMenuDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *AdminMenuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
