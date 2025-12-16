package admin

import (
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/samber/lo"

	"github.com/denghuo98/zzframe/internal/dao"
	"github.com/denghuo98/zzframe/internal/model/entity"
	"github.com/denghuo98/zzframe/web/zcontext"
	"github.com/denghuo98/zzframe/zconsts"
	"github.com/denghuo98/zzframe/zdb/zgorm"
	adminSchema "github.com/denghuo98/zzframe/zschema/admin"
	"github.com/denghuo98/zzframe/zservice"
)

type sAdminMember struct {
}

func init() {
	zservice.RegisterAdminMember(NewAdminMember())
}

func NewAdminMember() *sAdminMember {
	return &sAdminMember{}
}

// Info 获取当前登录用户信息
func (s *sAdminMember) Info(ctx g.Ctx) (out *adminSchema.LoginMemberInfoOutput, err error) {
	user := zcontext.GetUser(ctx)
	if user == nil {
		return nil, gerror.New("用户未登录")
	}
	out = &adminSchema.LoginMemberInfoOutput{
		Id:       user.Id,
		Username: user.Username,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Sex:      0, // 从上下文中无法获取性别，需要查询数据库
		Email:    user.Email,
		Mobile:   user.Mobile,
	}
	// 查询用户性别
	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(user.Id).Scan(&mb); err != nil {
		return nil, gerror.Wrap(err, zconsts.ErrorORM)
	}
	if mb != nil {
		out.Sex = mb.Sex
	}
	return out, nil
}

// List 获取用户列表
func (s *sAdminMember) List(ctx g.Ctx, in *adminSchema.MemberListInput) (out *adminSchema.MemberListOutput, total int, err error) {
	m := dao.AdminMember.Ctx(ctx)
	cols := dao.AdminMember.Columns()

	// 条件筛选
	if in.Username != "" {
		m = m.WhereLike(cols.Username, "%"+in.Username+"%")
	}
	if in.RealName != "" {
		m = m.WhereLike(cols.RealName, "%"+in.RealName+"%")
	}
	if in.Email != "" {
		m = m.WhereLike(cols.Email, "%"+in.Email+"%")
	}
	if in.Mobile != "" {
		m = m.WhereLike(cols.Mobile, "%"+in.Mobile+"%")
	}
	if in.Status > 0 {
		m = m.Where(cols.Status, in.Status)
	}

	// 查询总数
	count, err := m.Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, zconsts.ErrorORM)
	}
	total = count

	// 分页查询
	var list []*adminSchema.MemberListOutputItem
	if err = m.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return nil, 0, gerror.Wrap(err, zconsts.ErrorORM)
	}

	// 嵌入角色
	if err = s.embedRoles(ctx, list); err != nil {
		return nil, 0, err
	}

	out = &adminSchema.MemberListOutput{
		List: list,
	}
	return out, total, nil
}

func (s *sAdminMember) VerifySuperAdmin(ctx g.Ctx, id int64) bool {
	if id <= 0 {
		return false
	}
	var mb *entity.AdminMember
	if err := dao.AdminMember.Ctx(ctx).WherePri(id).Scan(&mb); err != nil {
		g.Log().Warning(ctx, "查询用户信息失败", err)
		return false
	}
	if mb == nil {
		return false
	}
	roleIds, err := dao.AdminMemberRole.Ctx(ctx).Where(dao.AdminMemberRole.Columns().MemberId, id).Array(dao.AdminMemberRole.Columns().RoleId)
	if err != nil {
		g.Log().Warning(ctx, "查询用户角色失败", err)
		return false
	}

	var roleEntities []*entity.AdminRole
	if err = dao.AdminRole.Ctx(ctx).WhereIn(dao.AdminRole.Columns().Id, roleIds).Scan(&roleEntities); err != nil {
		g.Log().Warning(ctx, "查询用户角色失败", err)
		return false
	}
	for _, role := range roleEntities {
		if role.Key == zconsts.SuperRoleKey {
			return true
		}
	}
	return false
}

func (s *sAdminMember) VerifyUnique(ctx g.Ctx, in *adminSchema.VerifyUniqueInput) (err error) {
	if in.Where == nil {
		return nil
	}

	cols := dao.AdminMember.Columns()
	msgMap := g.MapStrStr{
		cols.Username: "账号已存在，请更换一个",
		cols.Email:    "邮箱已存在，请更换一个",
		cols.Mobile:   "手机号码已存在，请更换一个",
	}

	for key, value := range in.Where {
		if value == "" {
			continue
		}
		message, ok := msgMap[key]
		if !ok {
			return gerror.Newf("字段 [%v] 未配置错误信息", key)
		}
		if err = zgorm.IsUnique(ctx, &dao.AdminMember, g.Map{key: value}, message, in.Id); err != nil {
			return err
		}
	}
	return nil
}

func (s *sAdminMember) Edit(ctx g.Ctx, in *adminSchema.MemberEditInput) (err error) {
	if in.Username == "" {
		return gerror.New("账号不能为空")
	}

	cols := dao.AdminMember.Columns()
	err = s.VerifyUnique(ctx, &adminSchema.VerifyUniqueInput{
		Id: in.Id,
		Where: g.Map{
			cols.Username: in.Username,
			cols.Email:    in.Email,
			cols.Mobile:   in.Mobile,
		},
	})
	if err != nil {
		return err
	}

	// 验证角色ID
	for _, roleId := range in.RoleIds {
		if err = zservice.AdminRole().VerifyRoleId(ctx, roleId); err != nil {
			return err
		}
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	m := g.Model(dao.AdminMember.Table()).TX(tx)
	// 修改
	if in.Id > 0 {
		// TODO 超管账号禁止编辑

		// 修改密码，需要使用到密码盐进行加密
		if in.Password != "" {
			salt, err := m.Fields(cols.Salt).WherePri(in.Id).Value()
			if err != nil {
				return gerror.Wrap(err, zconsts.ErrorORM)
			}
			if salt.IsEmpty() {
				return gerror.New("该用户没有设置密码盐，请联系管理员！")
			}
			in.PasswordHash = gmd5.MustEncryptString(in.Password + salt.String())
		} else {
			// 忽略密码字段
			m = m.FieldsEx(cols.PasswordHash)
		}

		if _, err = m.WherePri(in.Id).Data(in).Update(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		}
	} else {
		// 增加用户
		var addIn adminSchema.MemberAddInput
		addIn.MemberEditInput = in
		addIn.Salt = grand.S(6)
		addIn.PasswordHash = gmd5.MustEncryptString(in.Password + addIn.Salt)
		if id, err := m.Data(addIn).OmitEmpty().InsertAndGetId(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		} else {
			in.Id = id
		}
	}

	if len(in.RoleIds) > 0 {
		if err = s.updateRoles(tx, &adminSchema.MemberUpdateRoleInput{
			Id:      in.Id,
			RoleIds: in.RoleIds,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *sAdminMember) UpdateProfile(ctx g.Ctx, in *adminSchema.MemberUpdateProfileInput) (err error) {
	if in.Id <= 0 {
		return gerror.New("用户ID不能为空")
	}

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(in.Id).Scan(&mb); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if mb == nil {
		return gerror.New("用户不存在")
	}
	cols := dao.AdminMember.Columns()
	update := g.Map{
		cols.Avatar:   in.Avatar,
		cols.RealName: in.RealName,
		cols.Sex:      in.Sex,
	}
	if _, err = dao.AdminMember.Ctx(ctx).WherePri(in.Id).Data(update).Update(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	return nil
}

func (s *sAdminMember) UpdatePassword(ctx g.Ctx, in *adminSchema.MemberUpdatePasswordInput) (err error) {
	if in.Id <= 0 {
		return gerror.New("用户ID不能为空")
	}

	var mb *entity.AdminMember
	if err = dao.AdminMember.Ctx(ctx).WherePri(in.Id).Scan(&mb); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if mb == nil {
		return gerror.New("用户不存在")
	}

	if gmd5.MustEncryptString(in.OldPassword+mb.Salt) != mb.PasswordHash {
		return gerror.New("旧密码错误")
	}
	if gmd5.MustEncryptString(in.NewPassword+mb.Salt) == mb.PasswordHash {
		return gerror.New("新密码不能与旧密码相同")
	}

	update := g.Map{
		dao.AdminMember.Columns().PasswordHash: gmd5.MustEncryptString(in.NewPassword + mb.Salt),
	}
	if _, err = dao.AdminMember.Ctx(ctx).WherePri(in.Id).Data(update).Update(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	return nil
}

func (s *sAdminMember) Delete(ctx g.Ctx, in *adminSchema.MemberDeleteInput) (err error) {
	if in.Id <= 0 {
		return gerror.New("用户ID不能为空")
	}
	cols := dao.AdminMember.Columns()
	if ok, err := dao.AdminMember.Ctx(ctx).WherePri(in.Id).Exist(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	} else if !ok {
		return gerror.New("用户不存在")
	}
	// 软删除用户，将用户的状态设置成为禁用状态
	if _, err = dao.AdminMember.Ctx(ctx).WherePri(in.Id).Data(g.Map{cols.Status: zconsts.StatusDisable}).Update(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	return nil
}

func (s *sAdminMember) updateRoles(tx gdb.TX, in *adminSchema.MemberUpdateRoleInput) (err error) {
	// 校验用户是否存在
	if ok, err := g.Model(dao.AdminMember.Table()).TX(tx).Where(dao.AdminMember.Columns().Id, in.Id).Exist(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	} else if !ok {
		return gerror.New("用户不存在")
	}

	// 查询出来当前用户所绑定的角色ID
	mrCols := dao.AdminMemberRole.Columns()
	roleIds, err := g.Model(dao.AdminMemberRole.Table()).TX(tx).Where(mrCols.MemberId, in.Id).Array(mrCols.RoleId)
	if err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	existingRoleIds := make([]int64, 0, len(roleIds))
	for _, roleId := range roleIds {
		existingRoleIds = append(existingRoleIds, roleId.Int64())
	}
	// 将当前用户所绑定的角色ID与新角色ID进行对比，如果角色ID不存在，则删除
	needDeleteRoleIds, needInsertRoleIds := lo.Difference(existingRoleIds, in.RoleIds)

	if _, err = g.Model(dao.AdminMemberRole.Table()).TX(tx).Where(mrCols.MemberId, in.Id).WhereIn(mrCols.RoleId, needDeleteRoleIds).Delete(); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}
	if len(needInsertRoleIds) > 0 {
		data := make([]*entity.AdminMemberRole, 0, len(needInsertRoleIds))
		for _, roleId := range needInsertRoleIds {
			data = append(data, &entity.AdminMemberRole{
				MemberId: in.Id,
				RoleId:   roleId,
			})
		}
		if _, err = g.Model(dao.AdminMemberRole.Table()).TX(tx).Data(data).Insert(); err != nil {
			return gerror.Wrap(err, zconsts.ErrorORM)
		}
	}
	return nil
}

// embedRoles 查询中嵌入用户的角色列表
func (s *sAdminMember) embedRoles(ctx g.Ctx, list []*adminSchema.MemberListOutputItem) (err error) {
	memberIds := lo.Map(list, func(item *adminSchema.MemberListOutputItem, _ int) int64 {
		return item.Id
	})

	relations := make([]*entity.AdminMemberRole, 0)
	amCols := dao.AdminMemberRole.Columns()
	if err = dao.AdminMemberRole.Ctx(ctx).WhereIn(amCols.MemberId, memberIds).Scan(&relations); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}

	// 将数组转换成 Map 映射
	memberRoleIdsMap := lo.GroupByMap(relations, func(item *entity.AdminMemberRole) (int64, int64) {
		return item.MemberId, item.RoleId
	})

	// 查询对应的角色详情
	roleIds := lo.UniqMap(relations, func(item *entity.AdminMemberRole, _ int) int64 {
		return item.RoleId
	})

	roles := make([]*entity.AdminRole, 0, len(roleIds))
	if err = dao.AdminRole.Ctx(ctx).WhereIn(dao.AdminRole.Columns().Id, roleIds).Scan(&roles); err != nil {
		return gerror.Wrap(err, zconsts.ErrorORM)
	}

	// 绑定到对应的用户上
	for _, item := range list {
		roleIds := gconv.SliceInt64(memberRoleIdsMap[item.Id])
		roles := lo.Filter(roles, func(item *entity.AdminRole, _ int) bool {
			return lo.Contains(roleIds, item.Id)
		})
		item.RoleIds = roleIds
		item.Roles = roles
	}
	return
}
