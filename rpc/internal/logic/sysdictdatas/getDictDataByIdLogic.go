package sysdictdataslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictDataByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByIdLogic {
	return &GetDictDataByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictDataByIdLogic) GetDictDataById(in *apps.IdReq) (*apps.DictDataResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictDataResp{}, nil
}
