package systenantslogic

import (
	"context"
	"github.com/saas-zero/saas-zero-common/pkg/id"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/protobuf/proto"
)

type CreateTenantLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateTenantLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTenantLogic {
	return &CreateTenantLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateTenantLogic) CreateTenant(in *apps.TenantReq) (*apps.TenantResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	create := l.svcCtx.DB.SysTenant.Create().
		SetName(in.GetName()).
		SetCode(in.GetCode()).
		SetAdminID(in.GetAdminId()).
		SetStatus(systenant.Status(in.GetStatus()))

	if in.GetParentId() > 0 {
		create.SetParentID(in.GetParentId())
	}
	if in.GetPackageId() > 0 {
		create.SetPackageID(in.GetPackageId())
	}
	if in.GetExpiredAt() > 0 {
		create.SetExpiredAt(time.UnixMilli(in.GetExpiredAt()))
	}
	if in.GetRemark() != "" {
		create.SetRemark(in.GetRemark())
	}

	//logx.Infof("CreateTenant req:\n%s", prototext.Format(in))

	result, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.TenantResp{
		Code: int32(errno.Success.Code),
		Msg:  errno.Success.Msg,
		Data: &apps.Tenant{
			Id:    proto.Int64(result.ID),
			IdStr: proto.String(id.ToString(result.ID)),
		},
	}, nil
}
