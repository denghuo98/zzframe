// Package zservice 提供归档服务
package zservice

import (
	"context"
	"fmt"
	"time"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// ArchiveConfig 归档配置
type ArchiveConfig struct {
	// 归档阈值（天数），超过此天数的数据将被归档
	ArchiveDays int
	// 单次归档批量大小
	BatchSize int
}

// DefaultArchiveConfig 默认归档配置
var DefaultArchiveConfig = ArchiveConfig{
	ArchiveDays: 90, // 默认 90 天（3个月）
	BatchSize:   1000,
}

// ArchiveOperLog 归档操作日志
// 将指定天数之前的操作日志从主表移动到归档表
func ArchiveOperLog(ctx context.Context, config ...ArchiveConfig) (archivedCount int, err error) {
	cfg := DefaultArchiveConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	// 计算归档截止时间
	archiveBeforeTime := gtime.Now().Add(-time.Duration(cfg.ArchiveDays) * 24 * time.Hour)
	g.Log().Infof(ctx, "开始归档操作日志，归档 %s 之前的数据", archiveBeforeTime.Format("Y-m-d H:i:s"))

	// 开启事务
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		for {
			// 查询待归档的日志 ID
			idValues, err := tx.Model(dao.SysOperLog.Table()).
				Ctx(ctx).
				Fields("id").
				Where("created_at < ?", archiveBeforeTime.String()).
				Limit(cfg.BatchSize).
				Array()
			if err != nil {
				g.Log().Errorf(ctx, "查询待归档ID失败: %v", err)
				return err
			}

			if len(idValues) == 0 {
				g.Log().Infof(ctx, "没有更多需要归档的数据")
				break
			}

			// 转换为 int64 切片
			ids := make([]int64, len(idValues))
			for i, v := range idValues {
				ids[i] = v.Int64()
			}
			g.Log().Debugf(ctx, "待归档的ID: %v", ids)

			// 将数据插入归档表（查询后插入）
			// 将数据插入归档表（使用实体结构体确保字段映射正确）
			var logs []*entity.SysOperLog
			err = tx.Model(dao.SysOperLog.Table()).
				Ctx(ctx).
				WhereIn("id", ids).
				Scan(&logs)
			if err != nil {
				g.Log().Errorf(ctx, "查询待归档数据失败: %v", err)
				return err
			}
			g.Log().Debugf(ctx, "查询到 %d 条待归档记录", len(logs))

			if len(logs) > 0 {
				_, err = tx.Model("zz_sys_oper_log_archive").Ctx(ctx).Data(logs).Insert()
				if err != nil {
					g.Log().Errorf(ctx, "插入归档表失败: %v", err)
					return err
				}
				g.Log().Debugf(ctx, "成功插入归档表 %d 条", len(logs))
			}

			// 从主表删除已归档的数据
			result, err := tx.Model(dao.SysOperLog.Table()).
				Ctx(ctx).
				WhereIn("id", ids).
				Delete()
			if err != nil {
				g.Log().Errorf(ctx, "删除主表数据失败: %v", err)
				return err
			}
			deleted, _ := result.RowsAffected()
			g.Log().Debugf(ctx, "从主表删除 %d 条", deleted)

			archivedCount += len(ids)
			g.Log().Infof(ctx, "已归档 %d 条记录，累计 %d 条", len(ids), archivedCount)

			// 如果本次获取的数量小于批量大小，说明已经处理完毕
			if len(ids) < cfg.BatchSize {
				break
			}
		}
		return nil
	})

	if err != nil {
		g.Log().Errorf(ctx, "归档操作日志失败: %v", err)
		return 0, err
	}

	g.Log().Infof(ctx, "归档完成，共归档 %d 条操作日志", archivedCount)
	return archivedCount, nil
}

// ExportArchiveLog 导出归档日志到 JSON 文件
// 将归档表中指定时间之前的数据导出到文件
func ExportArchiveLog(ctx context.Context, beforeDays int, exportPath string) (exportedCount int, filePath string, err error) {
	if beforeDays <= 0 {
		beforeDays = 365 // 默认导出 1 年前的数据
	}
	if exportPath == "" {
		exportPath = "./data/archive"
	}

	// 确保导出目录存在
	if !gfile.Exists(exportPath) {
		if err = gfile.Mkdir(exportPath); err != nil {
			g.Log().Errorf(ctx, "创建导出目录失败: %v", err)
			return 0, "", err
		}
	}

	exportBeforeTime := gtime.Now().Add(-time.Duration(beforeDays) * 24 * time.Hour)
	g.Log().Infof(ctx, "开始导出归档日志，导出 %s 之前的数据", exportBeforeTime.Format("Y-m-d H:i:s"))

	// 查询待导出的数据
	var logs []*entity.SysOperLog
	err = g.DB().Model("zz_sys_oper_log_archive").
		Ctx(ctx).
		Where("archived_at < ?", exportBeforeTime.String()).
		Scan(&logs)
	if err != nil {
		g.Log().Errorf(ctx, "查询待导出数据失败: %v", err)
		return 0, "", err
	}

	if len(logs) == 0 {
		g.Log().Infof(ctx, "没有需要导出的数据")
		return 0, "", nil
	}

	// 生成文件名
	fileName := fmt.Sprintf("oper_log_%s.json", gtime.Now().Format("Ymd_His"))
	filePath = gfile.Join(exportPath, fileName)

	// 转换为 JSON
	jsonData, err := gjson.Encode(logs)
	if err != nil {
		g.Log().Errorf(ctx, "JSON 编码失败: %v", err)
		return 0, "", err
	}

	// 写入文件
	if err = gfile.PutBytes(filePath, jsonData); err != nil {
		g.Log().Errorf(ctx, "写入文件失败: %v", err)
		return 0, "", err
	}

	exportedCount = len(logs)
	g.Log().Infof(ctx, "导出完成，共导出 %d 条记录到 %s", exportedCount, filePath)
	return exportedCount, filePath, nil
}

// CleanArchiveLog 清理归档日志（冷数据处理）
// 先导出到文件，然后删除归档表中的数据
func CleanArchiveLog(ctx context.Context, retentionDays int, exportPath string) (exportedCount, deletedCount int, exportFile string, err error) {
	if retentionDays <= 0 {
		retentionDays = 365 // 默认保留 1 年
	}

	// 先导出数据
	exportedCount, exportFile, err = ExportArchiveLog(ctx, retentionDays, exportPath)
	if err != nil {
		return 0, 0, "", err
	}

	if exportedCount == 0 {
		g.Log().Infof(ctx, "没有需要清理的数据")
		return 0, 0, "", nil
	}

	// 删除已导出的数据
	cleanBeforeTime := gtime.Now().Add(-time.Duration(retentionDays) * 24 * time.Hour)
	g.Log().Infof(ctx, "开始清理归档日志，清理 %s 之前的数据", cleanBeforeTime.Format("Y-m-d H:i:s"))

	result, err := g.DB().Exec(ctx, `
		DELETE FROM zz_sys_oper_log_archive 
		WHERE archived_at < ?
	`, cleanBeforeTime.String())
	if err != nil {
		g.Log().Errorf(ctx, "清理归档日志失败: %v", err)
		return exportedCount, 0, exportFile, err
	}

	affected, _ := result.RowsAffected()
	deletedCount = int(affected)
	g.Log().Infof(ctx, "清理完成，导出 %d 条，删除 %d 条，文件: %s", exportedCount, deletedCount, exportFile)
	return exportedCount, deletedCount, exportFile, nil
}

// GetArchiveStats 获取归档统计信息
func GetArchiveStats(ctx context.Context) (stats map[string]interface{}, err error) {
	stats = make(map[string]interface{})

	// 主表记录数
	mainCount, err := g.DB().Model(dao.SysOperLog.Table()).Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	stats["main_count"] = mainCount

	// 归档表记录数
	archiveCount, err := g.DB().Model("zz_sys_oper_log_archive").Ctx(ctx).Count()
	if err != nil {
		return nil, err
	}
	stats["archive_count"] = archiveCount

	// 主表最早记录时间
	mainOldest, err := g.DB().Model(dao.SysOperLog.Table()).Ctx(ctx).Fields("MIN(created_at)").Value()
	if err == nil && !mainOldest.IsEmpty() {
		stats["main_oldest"] = mainOldest.String()
	}

	// 归档表最早记录时间
	archiveOldest, err := g.DB().Model("zz_sys_oper_log_archive").Ctx(ctx).Fields("MIN(created_at)").Value()
	if err == nil && !archiveOldest.IsEmpty() {
		stats["archive_oldest"] = archiveOldest.String()
	}

	return stats, nil
}
