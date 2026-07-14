package sysdeptslogic

import (
	"context"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/saas-zero/saas-zero-common/pkg/id"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type CreateDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDeptLogic {
	return &CreateDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDeptLogic) CreateDept(in *apps.DeptReq) (*apps.DeptResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)

	ctx := mixins.SetCurrentTenantId(l.ctx, tenantId)
	ctx = mixins.SetCurrentUserId(ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	create := l.svcCtx.DB.SysDept.Create().
		SetName(in.GetName()).
		SetStatus(sysdept.Status(in.GetStatus())).
		SetSort(uint32(in.GetSort())).
		SetMobile(in.GetMobile()).
		SetEmail(in.GetEmail())

	if in.GetParentId() > 0 {
		create.SetParentID(in.GetParentId())
	}
	if in.GetLeaderId() > 0 {
		create.SetLeaderID(in.GetLeaderId())
	}

	result, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.DeptResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: &apps.Dept{
			Id:    proto.Int64(result.ID),
			IdStr: proto.String(id.ToString(result.ID)),
		},
	}, nil
}
