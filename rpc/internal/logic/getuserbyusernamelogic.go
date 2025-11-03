package logic

import (
	"context"

	"system-service/rpc/apps/system-service"
	"system-service/rpc/internal/svc"

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

// 定义客户端流式 rpc
func (l *GetUserByUsernameLogic) GetUserByUsername(stream system_service.SysUsers_GetUserByUsernameServer) error {
	// todo: add your logic here and delete this line

	return nil
}
