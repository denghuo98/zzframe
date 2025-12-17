package zqueues

import (
	"context"
	"encoding/json"

	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zqueue"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/denghuo98/zzframe/zservice"
)

func init() {
	zqueue.RegisterConsumer(LoginLog)
}

// LoginLog 登录日志
var LoginLog = &qLoginLog{}

type qLoginLog struct{}

// GetTopic 主题
func (q *qLoginLog) GetTopic() string {
	return zconsts.QueueLoginLogTopic
}

// Handle 处理消息
func (q *qLoginLog) Handle(ctx context.Context, mqMsg zqueue.MqMsg) (err error) {
	var data entity.SysLoginLog
	if err = json.Unmarshal(mqMsg.Body, &data); err != nil {
		return err
	}
	return zservice.SysLoginLog().RealWrite(ctx, data)
}
