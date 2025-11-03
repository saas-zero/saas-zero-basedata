// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"system-service/api/internal/svc"
	"system-service/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type System_apiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSystem_apiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *System_apiLogic {
	return &System_apiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *System_apiLogic) System_api(req *types.Request) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
