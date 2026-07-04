package sysuserslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
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
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	u, err := l.svcCtx.DB.SysUser.TenantQuery(tenantId).
		Where(sysuser.UsernameEQ(in.GetUsername())).
		WithRoles().
		Only(l.ctx)
	if err != nil {
		return nil, err
	}
	return &apps.UserResp{
		Code: 200,
		Msg:  "success",
		Data: userToResp(u),
	}, nil
}
