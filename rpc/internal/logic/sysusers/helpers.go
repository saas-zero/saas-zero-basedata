package sysuserslogic

import (
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"google.golang.org/protobuf/proto"
)

func userToResp(u *ent.SysUser) *apps.User {
	resp := &apps.User{
		Id:          proto.Int64(u.ID),
		IdStr:       proto.String(strconv.FormatInt(u.ID, 10)),
		Username:    proto.String(u.Username),
		Nickname:    proto.String(u.Nickname),
		Mobile:      proto.String(u.Mobile),
		Email:       proto.String(u.Email),
		Status:      proto.String(string(u.Status)),
		LoginIp:     proto.String(u.LoginIP),
		TenantId:    proto.Int64(u.TenantID),
		TenantIdStr: proto.String(strconv.FormatInt(u.TenantID, 10)),
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
	if u.DeptID > 0 {
		resp.DeptId = proto.Int64(u.DeptID)
		resp.DeptIdStr = proto.String(strconv.FormatInt(u.DeptID, 10))
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
