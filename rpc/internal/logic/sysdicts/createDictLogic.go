package sysdictslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictLogic {
	return &CreateDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDictLogic) CreateDict(in *apps.DictReq) (*apps.DictResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictResp{}, nil
}
