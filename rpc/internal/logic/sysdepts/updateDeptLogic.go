package sysdeptslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdept"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
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
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	update := l.svcCtx.DB.SysDept.UpdateOneID(in.GetId())
	if in.Name != nil {
		update.SetName(in.GetName())
	}
	if in.ParentId != nil {
		if in.GetParentId() > 0 {
			update.SetParentID(in.GetParentId())
		}
	}
	if in.LeaderId != nil {
		if in.GetLeaderId() > 0 {
			update.SetLeaderID(in.GetLeaderId())
		}
	}
	if in.Mobile != nil {
		update.SetMobile(in.GetMobile())
	}
	if in.Email != nil {
		update.SetEmail(in.GetEmail())
	}
	if in.Status != nil {
		update.SetStatus(sysdept.Status(in.GetStatus()))
	}
	if in.Sort != nil {
		update.SetSort(uint32(in.GetSort()))
	}

	result, err := update.Save(ctx)
	if err != nil {
		return nil, err
	}

	d, _ := l.svcCtx.DB.SysDept.Query().Where(sysdept.IDEQ(result.ID)).WithLeader().Only(ctx)
	return &apps.DeptResp{
		Code: 200,
		Msg:  "success",
		Data: deptToResp(d),
	}, nil
}
