package sysdeptslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeptListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeptListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptListLogic {
	return &GetDeptListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeptListLogic) GetDeptList(in *apps.DeptPageReq) (*apps.DeptListResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	query := l.svcCtx.DB.SysDept.TenantQuery(tenantId)
	if in.GetName() != "" {
		query = query.Where(sysdept.NameContains(in.GetName()))
	}
	if in.GetStatus() != "" {
		query = query.Where(sysdept.StatusEQ(sysdept.Status(in.GetStatus())))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	page := int(in.GetPage())
	size := int(in.GetSize())
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	depts, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysdept.FieldSort)).
		WithLeader().
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.Dept, len(depts))
	for i, d := range depts {
		list[i] = deptToResp(d)
	}
	return &apps.DeptListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}
