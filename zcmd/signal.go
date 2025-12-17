package zcmd

import (
	"os"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"

	"github.com/denghuo98/zzframe/web/zutils"
	"github.com/denghuo98/zzframe/zconsts"
)

var (
	serverCloseSignal = make(chan struct{}, 1)
	serverWg          = sync.WaitGroup{}
	once              sync.Once
)

// signalHandlerForOverall 关闭信号处理
func SignalHandlerForOverall(sig os.Signal) {
	ServerCloseEvent(gctx.GetInitCtx())
	serverCloseSignal <- struct{}{}
}

// serverCloseEvent 关闭事件
// 区别于服务收到退出信号后的处理，只会执行一次
func ServerCloseEvent(ctx g.Ctx) {
	once.Do(func() {
		zutils.Event().Call(zconsts.EventServerClose, ctx)
	})
}

// SignalListen 信号监听
func SignalListen(ctx g.Ctx, handler ...gproc.SigHandler) {
	g.Go(ctx, func(ctx g.Ctx) {
		gproc.AddSigHandlerShutdown(handler...)
		gproc.Listen()
	}, func(ctx g.Ctx, exception error) {
		g.Log().Errorf(ctx, "signal listen failed, err:%+v", exception)
	})
}
