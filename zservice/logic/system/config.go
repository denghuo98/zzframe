package system

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmode"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/web/zcache"
	"github.com/denghuo98/zzframe/web/zstorager"
	"github.com/denghuo98/zzframe/web/ztoken"
	"github.com/denghuo98/zzframe/zconsts"
	webSchema "github.com/denghuo98/zzframe/zschema/zweb"
	"github.com/denghuo98/zzframe/zservice"
)

type sSystemConfig struct {
}

func init() {
	zservice.RegisterSystemConfig(NewSystemConfig())
}

func NewSystemConfig() *sSystemConfig {
	return &sSystemConfig{}
}

func (s *sSystemConfig) LoadConfig(ctx g.Ctx) (err error) {
	systemConfig := &webSchema.SystemConfig{}
	err = g.Cfg().MustGet(ctx, "system").Struct(systemConfig)
	if err != nil {
		return err
	}

	// 默认生产环境
	if systemConfig.Mode == "" {
		systemConfig.Mode = "prod"
	}
	s.SetGFMode(ctx, systemConfig.Mode)

	// 缓存配置
	cacheCfg, err := s.GetCacheConfig(ctx)
	if err != nil {
		return err
	}
	zcache.SetAdapter(ctx, cacheCfg)

	// token 认证配置
	tokenCfg, err := s.GetTokenConfig(ctx)
	if err != nil {
		return err
	}
	ztoken.SetConfig(tokenCfg)

	// 上传附件配置
	uploadCfg, err := s.GetUploadConfig(ctx)
	if err != nil {
		return err
	}
	zstorager.SetConfig(uploadCfg)

	// 初始化超级管理员
	if err := zservice.CommonSite().InitSuperAdmin(ctx); err != nil {
		return err
	}

	return nil
}

func (s *sSystemConfig) GetSuperAdmin(ctx g.Ctx) (conf *webSchema.SuperAdminConfig, err error) {
	conf = &webSchema.SuperAdminConfig{}
	err = g.Cfg().MustGet(ctx, "system.superAdmin").Struct(conf)
	return
}

func (s *sSystemConfig) GetCacheConfig(ctx g.Ctx) (conf *webSchema.CacheConfig, err error) {
	conf = &webSchema.CacheConfig{}
	v := g.Cfg().MustGet(ctx, "system.cache")
	if v != nil {
		err = v.Struct(conf)
	} else {
		conf = &webSchema.CacheConfig{
			Adapter: "file",
			FileDir: "tmp/cache",
		}
	}
	return
}

func (s *sSystemConfig) GetUploadConfig(ctx g.Ctx) (conf *webSchema.UploadConfig, err error) {
	conf = &webSchema.UploadConfig{}
	v := g.Cfg().MustGet(ctx, "system.upload")
	if v != nil {
		err = v.Struct(conf)
	} else {
		conf = &webSchema.UploadConfig{
			Drive:     zconsts.UploadDriveLocal,
			FileSize:  50,
			FileType:  "doc,docx,dot,xls,xlsx,xltx,ppt,pptx,pdf,txt,csv,html,xml,pptm,md",
			ImageSize: 10,
			ImageType: "jpg,png,jpeg,gif,webp",
			LocalPath: "upload-files-",
		}
	}
	return
}

func (s *sSystemConfig) GetTokenConfig(ctx g.Ctx) (conf *webSchema.TokenConfig, err error) {
	conf = &webSchema.TokenConfig{}
	v := g.Cfg().MustGet(ctx, "system.token")
	if v != nil {
		err = v.Struct(conf)
	} else {
		conf = &webSchema.TokenConfig{
			Expires:         86400,
			RefreshInterval: 3600,
			MaxRefreshTimes: 10,
			SecretKey:       "zzframe",
			MultiLogin:      true,
		}
	}
	return
}

func (s *sSystemConfig) SetGFMode(ctx g.Ctx, mode string) {
	if len(mode) == 0 {
		mode = gmode.NOT_SET
	}

	var modes = []string{gmode.DEVELOP, gmode.TESTING, gmode.STAGING, gmode.PRODUCT}

	if lo.Contains(modes, mode) {
		gmode.Set(mode)
	}
}
