package sysdictdataslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictDataLogic {
	return &UpdateDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDictDataLogic) UpdateDictData(in *apps.DictDataReq) (*apps.DictDataResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictDataResp{}, nil
}
