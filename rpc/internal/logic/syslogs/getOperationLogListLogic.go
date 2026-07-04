package syslogslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperationLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperationLogListLogic {
	return &GetOperationLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOperationLogListLogic) GetOperationLogList(in *apps.LogPageReq) (*apps.OperationLogListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.OperationLogListResp{}, nil
}
