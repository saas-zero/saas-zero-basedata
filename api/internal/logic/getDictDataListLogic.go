// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/api/internal/svc"
	"github.com/saas-zero/saas-zero-basedata/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDictDataListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataListLogic {
	return &GetDictDataListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDictDataListLogic) GetDictDataList(req *types.DictDataPageReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
