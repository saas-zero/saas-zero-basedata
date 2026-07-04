package sysdictslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictListLogic {
	return &GetDictListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictListLogic) GetDictList(in *apps.DictPageReq) (*apps.DictListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictListResp{}, nil
}
