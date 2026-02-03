// Package zqueues 提供队列消费者实现
package zqueues

import (
	"context"
	"encoding/json"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zqueue"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	zqueue.RegisterConsumer(OperLog)
}

// OperLog 操作日志消费者实例
var OperLog = &qOperLog{}

type qOperLog struct{}

// GetTopic 获取消费主题
func (q *qOperLog) GetTopic() string {
	return zconsts.QueueOperLogTopic
}

// Handle 处理队列消息
// 从队列获取操作日志数据，写入数据库
func (q *qOperLog) Handle(ctx context.Context, mqMsg zqueue.MqMsg) (err error) {
	var log entity.SysOperLog
	if err = json.Unmarshal(mqMsg.Body, &log); err != nil {
		g.Log().Errorf(ctx, "操作日志反序列化失败: %v", err)
		return err
	}

	_, err = dao.SysOperLog.Ctx(ctx).FieldsEx("id").Insert(log)
	if err != nil {
		g.Log().Errorf(ctx, "操作日志写入失败: %v", err)
	}
	return err
}
