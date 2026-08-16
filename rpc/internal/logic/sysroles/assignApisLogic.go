package sysroleslogic

import (
	"context"
	"fmt"
	"github.com/saas-zero/saas-zero-common/pkg/id"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"

	"github.com/saas-zero/saas-zero-basedata/ent/sysapi"
	"github.com/saas-zero/saas-zero-basedata/ent/sysrole"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
)

type AssignApisLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignApisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignApisLogic {
	return &AssignApisLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignApisLogic) AssignApis(in *apps.RoleReq) (*apps.EmptyResp, error) {
	roleCode := in.GetCode()
	tenantId := mixins.GetCurrentTenantId(l.ctx)
	dom := id.ToString(tenantId)

	// 前端只传 id + apiIds，不传 code，需从角色表查询兜底，
	// 否则会生成 v0 为空的脏策略。
	if roleCode == "" {
		role, err := l.svcCtx.DB.SysRole.Get(l.ctx, in.GetId())
		if err != nil {
			return nil, err
		}
		roleCode = role.Code
	}

	// API assignment is an edit of the role even though the association is
	// stored in Casbin. Touch the role so its audit fields stay authoritative.
	if err := l.svcCtx.DB.SysRole.UpdateOneID(in.GetId()).Exec(l.ctx); err != nil {
		return nil, err
	}

	l.svcCtx.Enforcer.RemoveFilteredPolicy(0, roleCode, dom)

	for _, apiId := range in.GetApiIds() {
		api, err := l.svcCtx.DB.SysApi.Get(l.ctx, apiId)
		if err != nil {
			continue
		}
		// 只对具体接口（api 类型）生成策略；目录（group）仅用于分组展示，不参与权限匹配
		if api.APIType != sysapi.APITypeAPI {
			continue
		}
		l.svcCtx.Enforcer.AddPolicy(roleCode, dom, api.APIPath, strings.ToUpper(string(api.APIMethod)), id.ToString(apiId))
	}

	users, err := l.svcCtx.DB.SysUser.Query().Where(sysuser.HasRolesWith(sysrole.CodeEQ(roleCode))).All(l.ctx)
	if err == nil {
		for _, u := range users {
			l.svcCtx.Redis.Incr(fmt.Sprintf("token_version:%d", u.ID))
		}
	}

	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
