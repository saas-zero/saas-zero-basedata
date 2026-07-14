package systenantslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/saas-zero/saas-zero-common/pkg/pagination"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetTenantListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTenantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTenantListLogic {
	return &GetTenantListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTenantListLogic) GetTenantList(in *apps.TenantPageReq) (*apps.TenantListResp, error) {
	query := l.svcCtx.DB.SysTenant.ActiveQuery()
	if in.GetName() != "" {
		query = query.Where(systenant.NameContains(in.GetName()))
	}
	if in.GetCode() != "" {
		query = query.Where(systenant.CodeContains(in.GetCode()))
	}
	if in.GetStatus() != "" {
		query = query.Where(systenant.StatusEQ(systenant.Status(in.GetStatus())))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	_, size, offset := pagination.Normalize(int(in.GetPage()), int(in.GetSize()))

	tenants, err := query.
		Offset(offset).
		Limit(size).
		Order(ent.Asc(systenant.FieldCreatedAt)).
		WithSysPackage().
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.Tenant, len(tenants))
	for i, t := range tenants {
		list[i] = tenantToResp(t)
	}
	return &apps.TenantListResp{
		Code:  int32(errno.Success.Code),
		Msg:   errno.Success.Msg,
		List:  list,
		Total: int64(total),
	}, nil
}
