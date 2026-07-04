package sysdictslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictByIdLogic {
	return &GetDictByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictByIdLogic) GetDictById(in *apps.IdReq) (*apps.DictResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictResp{}, nil
}
