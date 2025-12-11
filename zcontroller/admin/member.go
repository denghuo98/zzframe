package admin

import (
	"github.com/gogf/gf/v2/frame/g"

	adminApi "github.com/denghuo98/zzframe/zapi/admin"
	"github.com/denghuo98/zzframe/zservice"
)

var Member = cMember{}

type cMember struct{}

// Info 获取登录用户信息
func (c *cMember) Info(ctx g.Ctx, req *adminApi.MemberInfoReq) (res *adminApi.MemberInfoRes, err error) {
	out, err := zservice.AdminMember().Info(ctx)
	if err != nil {
		return nil, err
	}
	res = new(adminApi.MemberInfoRes)
	res.LoginMemberInfoOutput = out
	return res, nil
}

// List 获取用户列表
func (c *cMember) List(ctx g.Ctx, req *adminApi.MemberListReq) (res *adminApi.MemberListRes, err error) {
	list, totalCount, err := zservice.AdminMember().List(ctx, &req.MemberListInput)
	if err != nil {
		return nil, err
	}
	res = new(adminApi.MemberListRes)
	res.MemberListOutput = *list
	res.PageRes.Pack(req, int(totalCount))
	return res, nil
}

// Edit 修改/新增用户
func (c *cMember) Edit(ctx g.Ctx, req *adminApi.MemberEditReq) (res *adminApi.MemberEditRes, err error) {
	if err = zservice.AdminMember().Edit(ctx, &req.MemberEditInput); err != nil {
		return nil, err
	}
	return &adminApi.MemberEditRes{}, nil
}

// Delete 删除用户
func (c *cMember) Delete(ctx g.Ctx, req *adminApi.DeleteReq) (res *adminApi.DeleteRes, err error) {
	if err = zservice.AdminMember().Delete(ctx, &req.MemberDeleteInput); err != nil {
		return nil, err
	}
	return &adminApi.DeleteRes{}, nil
}

// UpdateProfile 更新用户资料
func (c *cMember) UpdateProfile(ctx g.Ctx, req *adminApi.UpdateProfileReq) (res *adminApi.UpdateProfileRes, err error) {
	if err = zservice.AdminMember().UpdateProfile(ctx, &req.MemberUpdateProfileInput); err != nil {
		return nil, err
	}
	return &adminApi.UpdateProfileRes{}, nil
}

// UpdatePwd 更新用户密码
func (c *cMember) UpdatePwd(ctx g.Ctx, req *adminApi.UpdatePwdReq) (res *adminApi.UpdatePwdRes, err error) {
	if err = zservice.AdminMember().UpdatePassword(ctx, &req.MemberUpdatePasswordInput); err != nil {
		return nil, err
	}
	return &adminApi.UpdatePwdRes{}, nil
}
