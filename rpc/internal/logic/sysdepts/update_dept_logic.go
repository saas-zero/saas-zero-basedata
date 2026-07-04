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

type UpdateDeptLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDeptLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDeptLogic {
	return &UpdateDeptLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDeptLogic) UpdateDept(in *apps.DeptReq) (*apps.DeptResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.DeptResp{Code: 500, Msg: "部门ID不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)

	update := l.svcCtx.DB.SysDept.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.ParentId != nil {
		update.SetParentID(*in.ParentId)
	}
	if in.LeaderId != nil {
		update.SetLeaderID(*in.LeaderId)
	}
	if in.Mobile != nil {
		update.SetMobile(*in.Mobile)
	}
	if in.Email != nil {
		update.SetEmail(*in.Email)
	}
	if in.Status != nil {
		update.SetStatus(sysdept.Status(*in.Status))
	}
	if in.Sort != nil {
		update.SetSort(uint32(*in.Sort))
	}

	dept, err := update.Save(l.ctx)
	if err != nil {
		return &apps.DeptResp{Code: 500, Msg: fmt.Sprintf("更新部门失败: %v", err)}, nil
	}

	return &apps.DeptResp{
		Code: 200,
		Msg:  "更新成功",
		Data: deptToPb(dept, tenantId),
	}, nil
}