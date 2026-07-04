package systenantslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

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

	if in.Name != nil && *in.Name != "" {
		query.Where(systenant.NameContains(*in.Name))
	}
	if in.Code != nil && *in.Code != "" {
		query.Where(systenant.CodeContains(*in.Code))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(systenant.StatusEQ(systenant.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.TenantListResp{Code: 500, Msg: fmt.Sprintf("查询租户总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	tenants, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Desc(systenant.FieldCreatedAt)).
		All(l.ctx)
	if err != nil {
		return &apps.TenantListResp{Code: 500, Msg: fmt.Sprintf("查询租户列表失败: %v", err)}, nil
	}

	list := make([]*apps.Tenant, 0, len(tenants))
	for _, t := range tenants {
		list = append(list, tenantToPb(t))
	}

	return &apps.TenantListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}