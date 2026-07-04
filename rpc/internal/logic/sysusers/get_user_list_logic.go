package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserListLogic) GetUserList(in *apps.UserPageReq) (*apps.UserListResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	query := l.svcCtx.DB.SysUser.TenantQuery(tenantId).
		WithSysDept().
		WithRoles()

	if in.Username != nil && *in.Username != "" {
		query.Where(sysuser.UsernameContains(*in.Username))
	}
	if in.Nickname != nil && *in.Nickname != "" {
		query.Where(sysuser.NicknameContains(*in.Nickname))
	}
	if in.Mobile != nil && *in.Mobile != "" {
		query.Where(sysuser.MobileContains(*in.Mobile))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysuser.StatusEQ(sysuser.Status(*in.Status)))
	}
	if in.DeptId != nil && *in.DeptId > 0 {
		query.Where(sysuser.DeptIDEQ(*in.DeptId))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.UserListResp{Code: 500, Msg: fmt.Sprintf("查询用户总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	users, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Desc(sysuser.FieldCreatedAt)).
		All(l.ctx)
	if err != nil {
		return &apps.UserListResp{Code: 500, Msg: fmt.Sprintf("查询用户列表失败: %v", err)}, nil
	}

	list := make([]*apps.User, 0, len(users))
	for _, u := range users {
		list = append(list, userToPb(u))
	}

	return &apps.UserListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}