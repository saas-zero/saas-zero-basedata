package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByUsernameLogic {
	return &GetUserByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByUsernameLogic) GetUserByUsername(in *apps.UserReq) (*apps.UserResp, error) {
	if in.Username == nil || *in.Username == "" {
		return &apps.UserResp{Code: 500, Msg: "用户名不能为空"}, nil
	}

	user, err := l.svcCtx.DB.SysUser.Query().
		Where(sysuser.UsernameEQ(*in.Username)).
		WithRoles().
		Only(l.ctx)
	if err != nil {
		return &apps.UserResp{Code: 500, Msg: fmt.Sprintf("获取用户失败: %v", err)}, nil
	}

	return &apps.UserResp{
		Code: 200,
		Msg:  "success",
		Data: userToPb(user),
	}, nil
}