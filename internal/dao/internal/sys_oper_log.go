// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysOperLogDao is the data access object for the table zz_sys_oper_log.
type SysOperLogDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysOperLogColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysOperLogColumns defines and stores column names for the table zz_sys_oper_log.
type SysOperLogColumns struct {
	Id         string // 日志ID
	Title      string // 操作标题
	Method     string // 请求方法
	RequestUri string // 请求URI
	OperParam  string // 请求参数
	Result     string // 响应结果
	OperId     string // 操作人ID
	OperName   string // 操作人名称
	OperIp     string // 操作IP
	BizId      string // 业务ID
	BizType    string // 业务类型
	Diff       string // 数据变更差异
	CreatedAt  string // 创建时间
	UpdatedAt  string // 修改时间
}

// sysOperLogColumns holds the columns for the table zz_sys_oper_log.
var sysOperLogColumns = SysOperLogColumns{
	Id:         "id",
	Title:      "title",
	Method:     "method",
	RequestUri: "request_uri",
	OperParam:  "oper_param",
	Result:     "result",
	OperId:     "oper_id",
	OperName:   "oper_name",
	OperIp:     "oper_ip",
	BizId:      "biz_id",
	BizType:    "biz_type",
	Diff:       "diff",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
}

// NewSysOperLogDao creates and returns a new DAO object for table data access.
func NewSysOperLogDao(handlers ...gdb.ModelHandler) *SysOperLogDao {
	return &SysOperLogDao{
		group:    "default",
		table:    "zz_sys_oper_log",
		columns:  sysOperLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysOperLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysOperLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysOperLogDao) Columns() SysOperLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysOperLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysOperLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SysOperLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
