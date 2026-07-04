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

type GetDeptTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeptTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptTreeLogic {
	return &GetDeptTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeptTreeLogic) GetDeptTree(in *apps.EmptyReq) (*apps.DeptTreeResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	depts, err := l.svcCtx.DB.SysDept.TenantQuery(tenantId).
		WithLeader().
		Order(ent.Asc(sysdept.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.DeptTreeResp{Code: 500, Msg: fmt.Sprintf("获取部门列表失败: %v", err)}, nil
	}

	tree := buildDeptTree(depts, 0, tenantId)

	return &apps.DeptTreeResp{
		Code: 200,
		Msg:  "success",
		Data: tree,
	}, nil
}

func buildDeptTree(depts []*ent.SysDept, parentId int64, tenantId int64) []*apps.Dept {
	tree := make([]*apps.Dept, 0)
	for _, d := range depts {
		if d.ParentID == parentId {
			dept := deptToPb(d, tenantId)
			dept.Children = buildDeptTree(depts, d.ID, tenantId)
			tree = append(tree, dept)
		}
	}
	return tree
}