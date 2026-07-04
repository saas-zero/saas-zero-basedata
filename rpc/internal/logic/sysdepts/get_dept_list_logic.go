package sysdeptslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

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

	if in.Name != nil && *in.Name != "" {
		query.Where(sysdept.NameContains(*in.Name))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysdept.StatusEQ(sysdept.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.DeptListResp{Code: 500, Msg: fmt.Sprintf("查询部门总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	depts, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysdept.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.DeptListResp{Code: 500, Msg: fmt.Sprintf("查询部门列表失败: %v", err)}, nil
	}

	list := make([]*apps.Dept, 0, len(depts))
	for _, d := range depts {
		list = append(list, deptToPb(d, tenantId))
	}

	return &apps.DeptListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}