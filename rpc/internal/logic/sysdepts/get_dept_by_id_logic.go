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

type GetDeptByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeptByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeptByIdLogic {
	return &GetDeptByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeptByIdLogic) GetDeptById(in *apps.IdReq) (*apps.DeptResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	dept, err := l.svcCtx.DB.SysDept.Query().
		Where(sysdept.IDEQ(int64(in.Id))).
		WithLeader().
		Only(l.ctx)
	if err != nil {
		return &apps.DeptResp{Code: 500, Msg: fmt.Sprintf("获取部门失败: %v", err)}, nil
	}

	return &apps.DeptResp{
		Code: 200,
		Msg:  "success",
		Data: deptToPb(dept, tenantId),
	}, nil
}

func deptToPb(d *ent.SysDept, tenantId int64) *apps.Dept {
	sort := int32(d.Sort)
	dept := &apps.Dept{
		Id:         &d.ID,
		IdStr:      strPtr(fmt.Sprintf("%d", d.ID)),
		Name:       &d.Name,
		ParentId:   &d.ParentID,
		ParentIdStr: strPtr(fmt.Sprintf("%d", d.ParentID)),
		LeaderId:   &d.LeaderID,
		LeaderIdStr: strPtr(fmt.Sprintf("%d", d.LeaderID)),
		Mobile:     &d.Mobile,
		Email:      &d.Email,
		Status:     strPtr(string(d.Status)),
		Sort:       &sort,
		TenantId:   &tenantId,
		TenantIdStr: strPtr(fmt.Sprintf("%d", tenantId)),
	}

	if d.Edges.Leader != nil {
		dept.LeaderName = &d.Edges.Leader.Nickname
	}

	createdAt := d.CreatedAt.Unix()
	dept.CreatedAt = &createdAt
	dept.CreatedBy = &d.CreatedBy

	updatedAt := d.UpdatedAt.Unix()
	dept.UpdatedAt = &updatedAt
	dept.UpdatedBy = &d.UpdatedBy

	return dept
}

func strPtr(s string) *string {
	return &s
}