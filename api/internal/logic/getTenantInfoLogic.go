package logic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/api/internal/svc"
	"github.com/saas-zero/saas-zero-basedata/api/internal/types"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTenantInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantInfoLogic {
	return &GetTenantInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTenantInfoLogic) GetTenantInfo(req *types.IdReq) (*types.BaseResp, error) {
	resp, err := l.svcCtx.SysTenants.GetTenantById(l.ctx, &apps.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{Code: int(resp.Code), Msg: resp.Msg, Data: resp.GetData()}, nil
}
