package logic

import (
	"context"
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/api/internal/svc"
	"github.com/saas-zero/saas-zero-basedata/api/internal/types"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type CreateTenantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTenantLogic) CreateTenant(req *types.TenantReq) (*types.BaseResp, error) {
	rpcReq := &apps.TenantReq{
		Name:   proto.String(req.Name),
		Code:   proto.String(req.Code),
		Status: proto.String(req.Status),
	}
	if req.AdminId > 0 {
		rpcReq.AdminId = proto.Int64(req.AdminId)
	}
	if req.ParentId > 0 {
		rpcReq.ParentId = proto.Int64(req.ParentId)
	}
	if req.PackageId > 0 {
		rpcReq.PackageId = proto.Int64(req.PackageId)
	}
	if req.ExpiredAt != "" {
		if v, err := strconv.ParseInt(req.ExpiredAt, 10, 64); err == nil {
			rpcReq.ExpiredAt = proto.Int64(v)
		}
	}
	if req.Remark != "" {
		rpcReq.Remark = proto.String(req.Remark)
	}
	resp, err := l.svcCtx.SysTenants.CreateTenant(l.ctx, rpcReq)
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{Code: int(resp.Code), Msg: resp.Msg, Data: resp.GetData()}, nil
}
