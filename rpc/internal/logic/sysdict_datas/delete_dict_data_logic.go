package sysdict_dataslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictDataLogic {
	return &DeleteDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDictDataLogic) DeleteDictData(in *apps.IdsReq) (*apps.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return &apps.EmptyResp{Code: 500, Msg: "字典数据ID不能为空"}, nil
	}

	_, err := l.svcCtx.DB.SysDictData.Update().
		Where(sysdictdata.IDIn(in.Ids...)).
		SetDeletedAt(time.Now()).
		Save(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("删除字典数据失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "删除成功"}, nil
}