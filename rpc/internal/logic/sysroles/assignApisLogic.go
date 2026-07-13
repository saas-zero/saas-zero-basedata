package sysroleslogic

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"strings"

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
	dom := strconv.FormatInt(tenantId, 10)

	l.svcCtx.Enforcer.RemoveFilteredPolicy(0, roleCode, dom)

	for _, apiId := range in.GetApiIds() {
		api, err := l.svcCtx.DB.SysApi.Get(l.ctx, apiId)
		if err != nil {
			continue
		}
		l.svcCtx.Enforcer.AddPolicy(roleCode, dom, api.APIPath, strings.ToUpper(string(api.APIMethod)), strconv.FormatInt(apiId, 10))
	}

	users, err := l.svcCtx.DB.SysUser.Query().Where(sysuser.HasRolesWith(sysrole.CodeEQ(roleCode))).All(l.ctx)
	if err == nil {
		for _, u := range users {
			l.svcCtx.Redis.Incr(fmt.Sprintf("token_version:%d", u.ID))
		}
	}

	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
