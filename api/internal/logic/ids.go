package logic

import (
	"github.com/saas-zero/saas-zero-basedata/api/internal/types"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	idutil "github.com/saas-zero/saas-zero-common/pkg/id"
)

// parseId converts a string ID (as returned by the frontend via `idStr`) into
// the int64 expected by the gRPC layer. JS loses precision for int64 (> 2^53),
// so IDs must travel as strings over JSON and only be parsed here on the wire.
// Implemented in saas-zero-common/pkg/id; kept as a package-local alias so the
// existing logic call sites don't need to change.
func parseId(s string) int64 {
	return idutil.Parse(s)
}

// parseIds converts a slice of string IDs to int64.
func parseIds(ss []string) []int64 {
	return idutil.ParseStrings(ss)
}

// formatIds converts int64 IDs to strings for JSON responses, avoiding the
// JS int64 precision loss that would occur if we returned raw int64.
func formatIds(ids []int64) []string {
	return idutil.ToStrings(ids)
}

// toSysUser maps the gRPC User (int64 roleIds) to the HTTP response type with
// string roleIds, so the frontend receives lossless IDs.
func toSysUser(u *apps.User) *types.SysUser {
	if u == nil {
		return nil
	}
	return &types.SysUser{
		Id:          u.GetIdStr(),
		IdStr:       u.GetIdStr(),
		Username:    u.GetUsername(),
		Nickname:    u.GetNickname(),
		Mobile:      u.GetMobile(),
		Email:       u.GetEmail(),
		DeptId:      u.GetDeptIdStr(),
		DeptIdStr:   u.GetDeptIdStr(),
		DeptName:    u.GetDeptName(),
		Status:      u.GetStatus(),
		Remark:      u.GetRemark(),
		LoginIp:     u.GetLoginIp(),
		LastLoginAt: idutil.ToString(u.GetLoginAt()),
		RoleIds:     formatIds(u.GetRoleIds()),
		RoleCodes:   u.GetRoleCodes(),
		RoleNames:   u.GetRoleNames(),
		TenantId:    u.GetTenantIdStr(),
		TenantIdStr: u.GetTenantIdStr(),
		CreatedAt:   idutil.ToString(u.GetCreatedAt()),
		CreatedBy:   u.GetCreatedBy(),
		UpdatedAt:   idutil.ToString(u.GetUpdatedAt()),
		UpdatedBy:   u.GetUpdatedBy(),
	}
}

// toSysUserList maps a slice of gRPC users.
func toSysUserList(users []*apps.User) []*types.SysUser {
	out := make([]*types.SysUser, 0, len(users))
	for _, u := range users {
		out = append(out, toSysUser(u))
	}
	return out
}
