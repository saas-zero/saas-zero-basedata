package sysapislogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiListLogic {
	return &GetApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApiListLogic) GetApiList(in *apps.ApiPageReq) (*apps.ApiListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.ApiListResp{}, nil
}
