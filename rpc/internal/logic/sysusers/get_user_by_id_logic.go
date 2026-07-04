package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserByIdLogic) GetUserById(in *apps.IdReq) (*apps.UserResp, error) {
	user, err := l.svcCtx.DB.SysUser.Query().
		Where(sysuser.IDEQ(int64(in.Id))).
		WithSysDept().
		WithRoles().
		Only(l.ctx)
	if err != nil {
		return &apps.UserResp{Code: 500, Msg: fmt.Sprintf("获取用户失败: %v", err)}, nil
	}

	return &apps.UserResp{
		Code: 200,
		Msg:  "success",
		Data: userToPb(user),
	}, nil
}

func userToPb(u *ent.SysUser) *apps.User {
	user := &apps.User{
		Id:        &u.ID,
		IdStr:     strPtr(fmt.Sprintf("%d", u.ID)),
		Username:  &u.Username,
		Nickname:  &u.Nickname,
		Mobile:    &u.Mobile,
		Email:     &u.Email,
		DeptId:    &u.DeptID,
		DeptIdStr: strPtr(fmt.Sprintf("%d", u.DeptID)),
		Status:    strPtr(string(u.Status)),
		Remark:    &u.Remark,
		LoginIp:   &u.LoginIP,
	}

	if !u.LoginAt.IsZero() {
		loginAt := u.LoginAt.Unix()
		user.LoginAt = &loginAt
	}

	if u.Edges.SysDept != nil {
		user.DeptName = &u.Edges.SysDept.Name
	}

	if len(u.Edges.Roles) > 0 {
		roleIds := make([]int64, 0, len(u.Edges.Roles))
		roleCodes := make([]string, 0, len(u.Edges.Roles))
		roleNames := make([]string, 0, len(u.Edges.Roles))
		for _, r := range u.Edges.Roles {
			roleIds = append(roleIds, r.ID)
			roleCodes = append(roleCodes, r.Code)
			roleNames = append(roleNames, r.Name)
		}
		user.RoleIds = roleIds
		user.RoleCodes = roleCodes
		user.RoleNames = roleNames
	}

	if u.TenantID != 0 {
		user.TenantId = &u.TenantID
		user.TenantIdStr = strPtr(fmt.Sprintf("%d", u.TenantID))
	}

	createdAt := u.CreatedAt.Unix()
	user.CreatedAt = &createdAt
	user.CreatedBy = &u.CreatedBy

	updatedAt := u.UpdatedAt.Unix()
	user.UpdatedAt = &updatedAt
	user.UpdatedBy = &u.UpdatedBy

	return user
}

func strPtr(s string) *string {
	return &s
}