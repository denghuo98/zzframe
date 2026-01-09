package zmigrate

import (
	"fmt"
	"path/filepath"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/util/gconv"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
)

func init() {
	InitDB(gctx.GetInitCtx())
}

// InitDB 初始化数据库表
// 根据配置文件中的 db 设置来使用不同的类型进行初始化
// 如果没有配置，则默认使用 sqlite，sqlite 文件路径自动生成在当前项目根目录下
// 此函数会自动设置 GF 的数据库配置，并在初始化完成后返回
func InitDB(ctx g.Ctx) error {
	// 读取数据库配置

	if !gdb.IsConfigured() {
		g.Log().Warning(ctx, "未配置数据库连接，使用默认 SQLite 数据库")
		gdb.SetConfig(gdb.Config{
			"default": getDefaultSQLiteConfig(),
		})
	}

	dbConfig, err := gdb.GetConfigGroup("default")
	if err != nil {
		g.Log().Panicf(ctx, "获取数据库配置失败: %v", err)
	}

	if err := initTables(ctx, dbConfig[0].Type); err != nil {
		return fmt.Errorf("初始化数据库表失败: %w", err)
	}

	g.Log().Infof(ctx, "数据库表初始化成功")
	return nil
}

// getDefaultSQLiteConfig 获取默认的 SQLite 配置
func getDefaultSQLiteConfig() gdb.ConfigGroup {
	// 获取当前项目根目录
	rootPath := gfile.Pwd()
	// SQLite 文件路径
	dbPath := filepath.Join(rootPath, "zzframe.db")

	return gdb.ConfigGroup{
		gdb.ConfigNode{
			Link:   fmt.Sprintf("sqlite::@file(%s)", dbPath),
			Prefix: "zz_",
		},
	}
}

// initTables 初始化数据库表
func initTables(ctx g.Ctx, dbType string) error {
	// 初始化管理员相关表
	tables := map[string]string{
		"zz_admin_member":      getCreateTableSQL("admin_member", dbType),
		"zz_admin_member_role": getCreateTableSQL("admin_member_role", dbType),
		"zz_admin_menu":        getCreateTableSQL("admin_menu", dbType),
		"zz_admin_role":        getCreateTableSQL("admin_role", dbType),
		"zz_admin_role_menu":   getCreateTableSQL("admin_role_menu", dbType),
		"zz_sys_login_log":     getCreateTableSQL("sys_login_log", dbType),
		"zz_sys_attachment":    getCreateTableSQL("sys_attachment", dbType),
	}

	for tableName, sql := range tables {
		// 检查表是否已创建
		if tableExists(ctx, tableName) {
			continue
		}

		// 创建表
		if err := createTable(ctx, sql, tableName); err != nil {
			return err
		}
	}

	return nil
}

// getCreateTableSQL 根据数据库类型获取建表 SQL
func getCreateTableSQL(tableKey, dbType string) string {
	switch dbType {
	case "mysql":
		switch tableKey {
		case "admin_member":
			return createAdminMemberTableSQL
		case "admin_member_role":
			return createAdminMemberRoleTableSQL
		case "admin_menu":
			return createAdminMenuTableSql
		case "admin_role":
			return createAdminRoleTableSQL
		case "admin_role_menu":
			return createAdminRoleMenuTableSQL
		case "sys_login_log":
			return createSysLoginLogTableSQL
		case "sys_attachment":
			return createSysAttachmentTableSQL
		}
	case "sqlite":
		switch tableKey {
		case "admin_member":
			return createAdminMemberTableSQLSQLite
		case "admin_member_role":
			return createAdminMemberRoleTableSQLSQLite
		case "admin_menu":
			return createAdminMenuTableSqlSQLite
		case "admin_role":
			return createAdminRoleTableSQLSQLite
		case "admin_role_menu":
			return createAdminRoleMenuTableSQLSQLite
		case "sys_login_log":
			return createSysLoginLogTableSQLSQLite
		case "sys_attachment":
			return createSysAttachmentTableSQLSQLite
		}
	}
	return ""
}

// tableExists 检查表是否已存在
func tableExists(ctx g.Ctx, tableName string) bool {
	var err error
	var result gdb.Record
	var count int

	// 获取表明前缀和类型
	config, _ := gdb.GetConfigGroup("default")

	switch config[0].Type {
	case "mysql":
		query := "SELECT COUNT(*) as count FROM information_schema.tables WHERE table_name = ?"
		result, err = g.DB().GetOne(ctx, query, tableName)
	case "sqlite":
		query := "SELECT COUNT(*) as count FROM sqlite_master WHERE type='table' AND name=?"
		result, err = g.DB().GetOne(ctx, query, tableName)
	}

	if err != nil {
		g.Log().Warningf(ctx, "检查表 %s 是否存在失败: %v", tableName, err)
		return false
	}
	count = gconv.Int(result["count"])

	return count > 0
}

// createTable 创建表
func createTable(ctx g.Ctx, sql, tableName string) error {

	if _, err := g.DB().Exec(ctx, sql); err != nil {
		return fmt.Errorf("创建表 %s 失败: %w", tableName, err)
	}

	g.Log().Infof(ctx, "表 %s 创建成功", tableName)
	return nil
}
