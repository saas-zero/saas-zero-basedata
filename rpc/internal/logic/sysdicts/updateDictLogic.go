package sysdictslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictLogic {
	return &UpdateDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDictLogic) UpdateDict(in *apps.DictReq) (*apps.DictResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictResp{}, nil
}
