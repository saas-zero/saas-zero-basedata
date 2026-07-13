package sysdictdataslogic

import (
	"context"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type CreateDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictDataLogic {
	return &CreateDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDictDataLogic) CreateDictData(in *apps.DictDataReq) (*apps.DictDataResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)

	ctx := mixins.SetCurrentTenantId(l.ctx, tenantId)
	ctx = mixins.SetCurrentUserId(ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	create := l.svcCtx.DB.SysDictData.Create().
		SetDictID(in.GetDictId()).
		SetName(in.GetName()).
		SetKey(in.GetKey()).
		SetValue(in.GetValue()).
		SetStatus(sysdictdata.Status(in.GetStatus()))

	if in.GetRemark() != "" {
		create.SetRemark(in.GetRemark())
	}

	result, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.DictDataResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: &apps.DictData{
			Id:    proto.Int64(result.ID),
			IdStr: proto.String(strconv.FormatInt(result.ID, 10)),
		},
	}, nil
}
