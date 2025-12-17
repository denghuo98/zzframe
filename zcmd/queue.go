package zcmd

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"

	"github.com/denghuo98/zzframe/web/zqueue"

	_ "github.com/denghuo98/zzframe/zqueues"
)

var (
	Queue = &gcmd.Command{
		Name:        "queue",
		Brief:       "消息队列",
		Description: ``,
		Func: func(ctx g.Ctx, parser *gcmd.Parser) (err error) {
			// 服务日志处理
			// queue.Logger().SetHandlers(global.LoggingServeLogHandler)

			g.Go(ctx, func(ctx g.Ctx) {
				zqueue.Logger().Debug(ctx, "start queue consumer..")
				zqueue.StartConsumersListener(ctx)
				zqueue.Logger().Debug(ctx, "start queue consumer success..")
			}, func(ctx g.Ctx, exception error) {
				zqueue.Logger().Errorf(ctx, "start queue consumer failed, err:%+v", exception)
			})

			serverWg.Add(1)

			// 信号监听
			SignalListen(ctx, SignalHandlerForOverall)

			<-serverCloseSignal
			zqueue.Logger().Debug(ctx, "queue successfully closed ..")
			serverWg.Done()
			return
		},
	}
)
