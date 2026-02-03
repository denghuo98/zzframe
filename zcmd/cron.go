package zcmd

import (
	"context"

	"github.com/denghuo98/zzframe/zservice"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gcron"
)

var (
	// Cron 定时任务命令
	Cron = &gcmd.Command{
		Name:        "cron",
		Brief:       "start cron server",
		Description: "启动定时任务服务，包含操作日志归档等定时任务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			g.Log().Info(ctx, "启动定时任务服务...")

			// 注册归档任务：每月1号凌晨2点执行
			_, err = gcron.Add(ctx, "0 0 2 1 * *", func(ctx context.Context) {
				g.Log().Info(ctx, "执行操作日志归档任务...")
				count, err := zservice.ArchiveOperLog(ctx)
				if err != nil {
					g.Log().Errorf(ctx, "操作日志归档失败: %v", err)
				} else {
					g.Log().Infof(ctx, "操作日志归档完成，共归档 %d 条", count)
				}
			}, "archive_oper_log")
			if err != nil {
				g.Log().Errorf(ctx, "注册归档任务失败: %v", err)
				return err
			}
			g.Log().Info(ctx, "归档任务已注册：每月1号凌晨2点执行")

			// 注册清理任务：每月15号凌晨3点执行（清理1年前的归档数据，先导出再删除）
			_, err = gcron.Add(ctx, "0 0 3 15 * *", func(ctx context.Context) {
				g.Log().Info(ctx, "执行归档日志清理任务...")
				exported, deleted, file, err := zservice.CleanArchiveLog(ctx, 365, "./data/archive")
				if err != nil {
					g.Log().Errorf(ctx, "归档日志清理失败: %v", err)
				} else {
					g.Log().Infof(ctx, "归档日志清理完成，导出 %d 条，删除 %d 条，文件: %s", exported, deleted, file)
				}
			}, "clean_archive_log")
			if err != nil {
				g.Log().Errorf(ctx, "注册清理任务失败: %v", err)
				return err
			}
			g.Log().Info(ctx, "清理任务已注册：每月15号凌晨3点执行")

			// 启动定时任务
			gcron.Start("archive_oper_log")
			gcron.Start("clean_archive_log")

			g.Log().Info(ctx, "定时任务服务启动成功")

			// 等待信号
			SignalListen(ctx, SignalHandlerForOverall)

			return nil
		},
	}

	// Archive 手动执行归档命令
	Archive = &gcmd.Command{
		Name:        "archive",
		Brief:       "archive operation logs",
		Description: "手动执行操作日志归档任务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// 获取归档天数参数，默认 90 天
			days := parser.GetOpt("days", "90").Int()
			g.Log().Infof(ctx, "手动执行归档，归档 %d 天前的数据", days)

			config := zservice.ArchiveConfig{
				ArchiveDays: days,
				BatchSize:   1000,
			}
			count, err := zservice.ArchiveOperLog(ctx, config)
			if err != nil {
				g.Log().Errorf(ctx, "归档失败: %v", err)
				return err
			}
			g.Log().Infof(ctx, "归档完成，共归档 %d 条记录", count)
			return nil
		},
	}

	// ArchiveExport 导出归档数据到文件
	ArchiveExport = &gcmd.Command{
		Name:        "archive-export",
		Brief:       "export archived logs to file",
		Description: "将归档的操作日志导出到 JSON 文件",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			days := parser.GetOpt("days", "365").Int()
			path := parser.GetOpt("path", "./data/archive").String()
			g.Log().Infof(ctx, "导出 %d 天前的归档数据到 %s", days, path)

			count, file, err := zservice.ExportArchiveLog(ctx, days, path)
			if err != nil {
				g.Log().Errorf(ctx, "导出失败: %v", err)
				return err
			}
			if count == 0 {
				g.Log().Info(ctx, "没有需要导出的数据")
			} else {
				g.Log().Infof(ctx, "导出完成，共 %d 条，文件: %s", count, file)
			}
			return nil
		},
	}

	// ArchiveClean 清理归档数据（导出后删除）
	ArchiveClean = &gcmd.Command{
		Name:        "archive-clean",
		Brief:       "export and clean archived logs",
		Description: "将归档的操作日志导出后删除（冷数据处理）",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			days := parser.GetOpt("days", "365").Int()
			path := parser.GetOpt("path", "./data/archive").String()
			g.Log().Infof(ctx, "清理 %d 天前的归档数据，导出到 %s", days, path)

			exported, deleted, file, err := zservice.CleanArchiveLog(ctx, days, path)
			if err != nil {
				g.Log().Errorf(ctx, "清理失败: %v", err)
				return err
			}
			if exported == 0 {
				g.Log().Info(ctx, "没有需要清理的数据")
			} else {
				g.Log().Infof(ctx, "清理完成，导出 %d 条，删除 %d 条，文件: %s", exported, deleted, file)
			}
			return nil
		},
	}

	// ArchiveStats 查看归档统计
	ArchiveStats = &gcmd.Command{
		Name:        "archive-stats",
		Brief:       "show archive statistics",
		Description: "查看操作日志归档统计信息",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			stats, err := zservice.GetArchiveStats(ctx)
			if err != nil {
				g.Log().Errorf(ctx, "获取统计信息失败: %v", err)
				return err
			}
			g.Log().Infof(ctx, "归档统计信息:")
			g.Log().Infof(ctx, "  主表记录数: %v", stats["main_count"])
			g.Log().Infof(ctx, "  归档表记录数: %v", stats["archive_count"])
			if oldest, ok := stats["main_oldest"]; ok {
				g.Log().Infof(ctx, "  主表最早记录: %v", oldest)
			}
			if oldest, ok := stats["archive_oldest"]; ok {
				g.Log().Infof(ctx, "  归档表最早记录: %v", oldest)
			}
			return nil
		},
	}
)

