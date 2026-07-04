package sysdictslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdict"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDictLogic {
	return &DeleteDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteDictLogic) DeleteDict(in *apps.IdsReq) (*apps.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return &apps.EmptyResp{Code: 500, Msg: "字典ID不能为空"}, nil
	}

	_, err := l.svcCtx.DB.SysDict.Update().
		Where(sysdict.IDIn(in.Ids...)).
		SetDeletedAt(time.Now()).
		Save(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("删除字典失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "删除成功"}, nil
}