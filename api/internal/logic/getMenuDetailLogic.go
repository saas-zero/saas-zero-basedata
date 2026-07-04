// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/api/internal/svc"
	"github.com/saas-zero/saas-zero-basedata/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuDetailLogic {
	return &GetMenuDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuDetailLogic) GetMenuDetail(req *types.IdReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
