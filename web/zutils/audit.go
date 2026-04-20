// Package zutils 提供公共工具函数
package zutils

import (
	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/gogf/gf/v2/frame/g"
)

// AuditBatchLimit 批量审计上限
const AuditBatchLimit = 100

// AuditData 审计数据结构
type AuditData struct {
	BizId   string      `json:"bizId"`   // 业务ID
	BizType string      `json:"bizType"` // 业务类型
	Old     interface{} `json:"old"`     // 变更前数据
	New     interface{} `json:"new"`     // 变更后数据
}

// NewAudit 创建审计数据
//
// 参数:
//
//	bizType: 业务类型（如 "admin_member", "sys_config"）
//	bizId: 业务ID（如用户ID、配置键）
//	oldData: 变更前数据
//	newData: 变更后数据
func NewAudit(bizType, bizId string, oldData, newData interface{}) AuditData {
	return AuditData{
		BizType: bizType,
		BizId:   bizId,
		Old:     oldData,
		New:     newData,
	}
}

// SetAudit 设置审计数据到上下文
// 支持单条和批量操作，使用变长参数
// 批量操作最多记录 100 条
//
// 使用示例:
//
//	单条: SetAudit(ctx, NewAudit("user", "1001", old, new))
//	批量: SetAudit(ctx, item1, item2, item3...)
func SetAudit(ctx g.Ctx, items ...AuditData) {
	if len(items) == 0 {
		return
	}

	// 批量上限检查
	if len(items) > AuditBatchLimit {
		g.Log().Warningf(ctx, "审计数据超过上限%d条，已截断", AuditBatchLimit)
		items = items[:AuditBatchLimit]
	}

	zcontext.SetData(ctx, string(zconsts.ContextKeyAuditData), items)
}

// GetAudit 从上下文获取审计数据
func GetAudit(ctx g.Ctx) []AuditData {
	data := zcontext.GetData(ctx, string(zconsts.ContextKeyAuditData))
	if data == nil {
		return nil
	}

	// 从 Data map 中获取审计数据
	if auditData, ok := data[string(zconsts.ContextKeyAuditData)]; ok {
		if items, ok := auditData.([]AuditData); ok {
			return items
		}
	}

	return nil
}
