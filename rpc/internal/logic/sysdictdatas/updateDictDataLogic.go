package sysdictdataslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
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
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	update := l.svcCtx.DB.SysDictData.UpdateOneID(in.GetId())
	if in.Name != nil {
		update.SetName(in.GetName())
	}
	if in.Key != nil {
		update.SetKey(in.GetKey())
	}
	if in.Value != nil {
		update.SetValue(in.GetValue())
	}
	if in.Status != nil {
		update.SetStatus(sysdictdata.Status(in.GetStatus()))
	}
	if in.Remark != nil {
		update.SetRemark(in.GetRemark())
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	d, err := l.svcCtx.DB.SysDictData.Query().Where(sysdictdata.IDEQ(result.ID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.DictDataResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: dictDataToResp(d),
	}, nil
}
