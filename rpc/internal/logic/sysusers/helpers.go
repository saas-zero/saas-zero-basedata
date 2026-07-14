package sysuserslogic

import (
	"github.com/saas-zero/saas-zero-common/pkg/id"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"google.golang.org/protobuf/proto"
)

func userToResp(u *ent.SysUser) *apps.User {
	resp := &apps.User{
		Id:          proto.Int64(u.ID),
		IdStr:       proto.String(id.ToString(u.ID)),
		Username:    proto.String(u.Username),
		Nickname:    proto.String(u.Nickname),
		Mobile:      proto.String(u.Mobile),
		Email:       proto.String(u.Email),
		Status:      proto.String(string(u.Status)),
		LoginIp:     proto.String(u.LoginIP),
		TenantId:    proto.Int64(u.TenantID),
		TenantIdStr: proto.String(id.ToString(u.TenantID)),
		CreatedAt:   proto.Int64(u.CreatedAt.UnixMilli()),
		UpdatedAt:   proto.Int64(u.UpdatedAt.UnixMilli()),
	}
	if u.CreatedBy != "" {
		resp.CreatedBy = proto.String(u.CreatedBy)
	}
	if u.UpdatedBy != "" {
		resp.UpdatedBy = proto.String(u.UpdatedBy)
	}
	if u.Remark != "" {
		resp.Remark = proto.String(u.Remark)
	}
	// SECURITY: Password is intentionally excluded from generic gRPC responses.
	// The bcrypt hash should never leak through list or detail queries.
	// Only GetUserByUsername (login flow) sets Password separately after this function returns.
	// See: getUserByUsernameLogic.go
	if u.DeptID > 0 {
		resp.DeptId = proto.Int64(u.DeptID)
		resp.DeptIdStr = proto.String(id.ToString(u.DeptID))
	}
	if !u.LoginAt.IsZero() {
		resp.LoginAt = proto.Int64(u.LoginAt.UnixMilli())
	}
	if len(u.Edges.Roles) > 0 {
		roleIds := make([]int64, len(u.Edges.Roles))
		roleCodes := make([]string, len(u.Edges.Roles))
		roleNames := make([]string, len(u.Edges.Roles))
		for i, r := range u.Edges.Roles {
			roleIds[i] = r.ID
			roleCodes[i] = r.Code
			roleNames[i] = r.Name
		}
		resp.RoleIds = roleIds
		resp.RoleCodes = roleCodes
		resp.RoleNames = roleNames
	}
	return resp
}
