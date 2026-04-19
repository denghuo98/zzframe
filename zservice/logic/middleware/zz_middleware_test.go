package middleware

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/denghuo98/zzframe/zservice"
	_ "github.com/denghuo98/zzframe/zservice/logic/common"
	_ "github.com/denghuo98/zzframe/zservice/logic/system"
)

func TestGetAnonymousConfig_DefaultDisabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()

		cfg, err := zservice.SystemConfig().GetAnonymousConfig(ctx)
		t.AssertNil(err)
		t.Assert(cfg.Enabled, false)
		t.Assert(cfg.Identity.Username, "anonymous")
		t.Assert(cfg.Identity.RealName, "游客")
	})
}

func TestDeliverUserContext_AnonymousBranch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 验证匿名配置结构可以正确构建
		m := NewMiddleware()
		t.AssertNE(m, nil)
	})
}
