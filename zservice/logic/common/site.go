package common

import (
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/web/zencrypt"
	"github.com/denghuo98/zzframe/web/ztoken"
	"github.com/denghuo98/zzframe/zconsts"
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	commonSchema "github.com/denghuo98/zzframe/zschema/common"
	systemSchema "github.com/denghuo98/zzframe/zschema/system"
	"github.com/denghuo98/zzframe/zschema/zweb"
	"github.com/denghuo98/zzframe/zservice"
)

type sCommonSite struct {
}

func init() {
	zservice.RegisterCommonSite(NewCommonSite())
}

func NewCommonSite() *sCommonSite {
	return &sCommonSite{}
}

// Ping 心跳检测
func (s *sCommonSite) Ping(ctx g.Ctx) (status string) {
	return "Pong"
}

func (s *sCommonSite) InitSuperAdmin(ctx g.Ctx) (err error) {
	conf, err := zservice.SystemConfig().GetSuperAdmin(ctx)
	if err != nil {
		return err
	}
	// 先判断是否存在角色
	var roleEntity *entity.AdminRole
	if err = dao.AdminRole.Ctx(ctx).Where(dao.AdminRole.Columns().Key, zconsts.SuperRoleKey).Scan(&roleEntity); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if roleEntity == nil {
		roleEntity, err = zservice.AdminRole().Edit(ctx, adminSchema.RoleEditInput{
			AdminRole: entity.AdminRole{
				Name:   zconsts.SuperRoleKey,
				Key:    zconsts.SuperRoleKey,
				Remark: "系统超级管理员角色，禁止修改删除",
				Sort:   1,
				Status: 1,
			},
		})
		if err != nil {
			return err
		}
		g.Log().Info(ctx, "已自动创建超级管理员角色")
	}

	// 如果配置文件中没有设置密码，使用默认值
	if conf.Password == "" {
		conf.Password = zconsts.DefaultSuperAdminPassword
	}
	// 判断是否存在，不存在则创建
	var memberEntity *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Username, zconsts.DefaultSuperAdminUsername).Scan(&memberEntity); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if memberEntity == nil {
		payload := adminSchema.MemberEditInput{
			Username: zconsts.DefaultSuperAdminUsername,
			Password: conf.Password,
			RealName: "超级管理员",
			Email:    "",
			Sex:      3,
			RoleIds:  []int64{roleEntity.Id},
			Remark:   "系统超级管理员，禁止编辑或删除",
			Status:   1,
		}
		err = zservice.AdminMember().Edit(ctx, &payload)
		if err != nil {
			return err
		}
		g.Log().Info(ctx, "已自动创建超级管理员用户")
	} else {
		// 只更新密码
		passwordHash := gmd5.MustEncryptString(conf.Password + memberEntity.Salt)
		if passwordHash != memberEntity.PasswordHash {
			dao.AdminMember.Ctx(ctx).WherePri(memberEntity.Id).Data(g.Map{
				dao.AdminMember.Columns().PasswordHash: passwordHash,
			}).Update()
			g.Log().Info(ctx, "已自动更新超级管理员密码")
		}
	}

	return nil
}

func (s *sCommonSite) AccountLogin(ctx g.Ctx, in *commonSchema.SiteAccountLoginInput) (out *commonSchema.SiteLoginOutput, err error) {
	defer func() {
		// 推送登录事件
		zservice.SysLoginLog().Push(ctx, &systemSchema.SysLoginLogPushInput{
			Response: out,
			Error:    err,
		})
	}()

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Username, in.Username).Scan(&mb); err != nil {
		err = gerror.Wrap(err, zconsts.ErrorORM)
		return
	}

	if mb == nil {
		err = gerror.New("账号不存在")
		return
	}

	if mb.Status != zconsts.StatusEnabled {
		err = gerror.New("账号已禁用")
		return
	}

	out = new(commonSchema.SiteLoginOutput)
	out.Id = mb.Id
	out.Username = mb.Username
	if mb.Salt == "" {
		err = gerror.New("该用户没有设置密码盐，请联系管理员！")
		return
	}

	if err = s.verifyPassword(in.Password, mb.Salt, mb.PasswordHash); err != nil {
		return
	}

	return s.handleLogin(ctx, mb)
}

// verifyPassword 验证密码
func (s *sCommonSite) verifyPassword(input, salt, hash string) (err error) {
	// 先解密文本
	unBase64, err := gbase64.Decode([]byte(input))
	if err != nil {
		return err
	}

	unAes, err := zencrypt.AesECBDecrypt(unBase64, []byte(zconsts.RequestEncryptKey))
	if err != nil {
		return err
	}

	plainText := string(unAes)
	if hash != gmd5.MustEncryptString(plainText+salt) {
		return gerror.New("用户密码错误")
	}

	return nil
}

func (s *sCommonSite) handleLogin(ctx g.Ctx, mb *entity.AdminMember) (out *commonSchema.SiteLoginOutput, err error) {
	// 获取用户绑定的角色
	roles, err := s.getMemberRoles(ctx, mb.Id)
	if err != nil {
		return nil, err
	}

	identity := &zweb.Identity{
		Id:       mb.Id,
		Roles:    roles,
		Username: mb.Username,
		RealName: mb.RealName,
		Avatar:   mb.Avatar,
		Email:    mb.Email,
		Mobile:   mb.Mobile,
		App:      "app",
		LoginAt:  gtime.Now(),
	}

	tokenStr, expires, err := ztoken.Login(ctx, identity)
	if err != nil {
		return nil, err
	}
	out = new(commonSchema.SiteLoginOutput)
	out.Id = mb.Id
	out.Username = mb.Username
	out.Token = tokenStr
	out.Expires = expires

	return out, nil
}

func (s *sCommonSite) getMemberRoles(ctx g.Ctx, id int64) (roles []*zweb.IdentityRole, err error) {
	roleIds, err := dao.AdminMemberRole.Ctx(ctx).Where(dao.AdminMemberRole.Columns().MemberId, id).Array(dao.AdminMemberRole.Columns().RoleId)
	if err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}

	var roleEntities []*entity.AdminRole
	if err = dao.AdminRole.Ctx(ctx).WhereIn(dao.AdminRole.Columns().Id, roleIds).Scan(&roleEntities); err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}
	for _, role := range roleEntities {
		roles = append(roles, &zweb.IdentityRole{
			Id:   role.Id,
			Name: role.Name,
			Key:  role.Key,
		})
	}
	return roles, nil
}

// BindUserContext 将用户信息绑定到上下文中
// 主要是重新查询数据库获取用户信息，因为用户信息可能发生变化
func (s *sCommonSite) BindUserContext(ctx g.Ctx, claims *zweb.Identity) (err error) {
	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).Where(dao.AdminMember.Columns().Id, claims.Id).Scan(&mb); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if mb == nil {
		return gerror.New("用户不存在")
	}
	if mb.Status != zconsts.StatusEnabled {
		return gerror.New("用户已禁用，如有疑问请联系管理员")
	}

	roles, err := s.getMemberRoles(ctx, mb.Id)
	if err != nil {
		return err
	}

	user := &zweb.Identity{
		Id:       mb.Id,
		Username: mb.Username,
		Roles:    roles,
		RealName: mb.RealName,
		Avatar:   mb.Avatar,
		Email:    mb.Email,
		Mobile:   mb.Mobile,
		App:      claims.App,
		LoginAt:  claims.LoginAt,
	}

	zcontext.SetUser(ctx, user)

	return nil
}
