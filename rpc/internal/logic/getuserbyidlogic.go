package logic

import (
	"context"

	"system-service/rpc/apps/system-service"
	"system-service/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 定义客户端流式 rpc
func (l *GetUserByIdLogic) GetUserById(stream system_service.SysUsers_GetUserByIdServer) error {
	// todo: add your logic here and delete this line

	return nil
}
