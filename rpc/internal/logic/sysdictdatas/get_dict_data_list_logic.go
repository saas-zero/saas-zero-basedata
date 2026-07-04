package sysdictdataslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictDataListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataListLogic {
	return &GetDictDataListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictDataListLogic) GetDictDataList(in *apps.DictDataPageReq) (*apps.DictDataListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.DictDataListResp{}, nil
}
