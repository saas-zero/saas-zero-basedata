package sysdictdataslogic

import (
	"context"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
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
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	_, err := l.svcCtx.DB.SysDictData.Update().
		Where(sysdictdata.IDIn(in.GetIds()...)).
		SetDeletedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
