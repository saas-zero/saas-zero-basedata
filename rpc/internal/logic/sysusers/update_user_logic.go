package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserLogic {
	return &UpdateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserLogic) UpdateUser(in *apps.UserReq) (*apps.UserResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.UserResp{Code: 500, Msg: "用户ID不能为空"}, nil
	}

	update := l.svcCtx.DB.SysUser.UpdateOneID(int64(*in.Id))

	if in.Nickname != nil {
		update.SetNickname(*in.Nickname)
	}
	if in.Mobile != nil {
		update.SetMobile(*in.Mobile)
	}
	if in.Email != nil {
		update.SetEmail(*in.Email)
	}
	if in.DeptId != nil {
		update.SetDeptID(*in.DeptId)
	}
	if in.Status != nil {
		update.SetStatus(sysuser.Status(*in.Status))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	user, err := update.Save(l.ctx)
	if err != nil {
		return &apps.UserResp{Code: 500, Msg: fmt.Sprintf("更新用户失败: %v", err)}, nil
	}

	return &apps.UserResp{
		Code: 200,
		Msg:  "更新成功",
		Data: userToPb(user),
	}, nil
}