package sysdeptslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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
	if in.Name == nil || *in.Name == "" {
		return &apps.DeptResp{Code: 500, Msg: "部门名称不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)
	ctx := mixins.SetCurrentTenantId(l.ctx, tenantId)

	create := l.svcCtx.DB.SysDept.Create().
		SetName(*in.Name)

	if in.ParentId != nil {
		create.SetParentID(*in.ParentId)
	}
	if in.LeaderId != nil {
		create.SetLeaderID(*in.LeaderId)
	}
	if in.Mobile != nil {
		create.SetMobile(*in.Mobile)
	}
	if in.Email != nil {
		create.SetEmail(*in.Email)
	}
	if in.Status != nil {
		create.SetStatus(sysdept.Status(*in.Status))
	}
	if in.Sort != nil {
		create.SetSort(uint32(*in.Sort))
	}

	dept, err := create.Save(ctx)
	if err != nil {
		return &apps.DeptResp{Code: 500, Msg: fmt.Sprintf("创建部门失败: %v", err)}, nil
	}

	return &apps.DeptResp{
		Code: 200,
		Msg:  "创建成功",
		Data: deptToPb(dept, tenantId),
	}, nil
}