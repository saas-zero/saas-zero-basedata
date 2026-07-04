package sysdeptslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeptByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeptByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptByIdLogic {
	return &GetDeptByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeptByIdLogic) GetDeptById(in *apps.IdReq) (*apps.DeptResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DeptResp{}, nil
}
