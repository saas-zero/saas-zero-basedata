package sysroleslogic

import (
	"strconv"

	casbinapi "github.com/casbin/casbin/v2"
	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"google.golang.org/protobuf/proto"
)

func roleToResp(r *ent.SysRole) *apps.Role {
	resp := &apps.Role{
		Id:          proto.Int64(r.ID),
		IdStr:       proto.String(strconv.FormatInt(r.ID, 10)),
		Name:        proto.String(r.Name),
		Code:        proto.String(r.Code),
		Status:      proto.String(string(r.Status)),
		Sort:        proto.Int32(int32(r.Sort)),
		TenantId:    proto.Int64(r.TenantID),
		TenantIdStr: proto.String(strconv.FormatInt(r.TenantID, 10)),
		CreatedAt:   proto.Int64(r.CreatedAt.UnixMilli()),
		UpdatedAt:   proto.Int64(r.UpdatedAt.UnixMilli()),
	}
	if r.Remark != "" {
		resp.Remark = proto.String(r.Remark)
	}
	if r.CreatedBy != "" {
		resp.CreatedBy = proto.String(r.CreatedBy)
	}
	if r.UpdatedBy != "" {
		resp.UpdatedBy = proto.String(r.UpdatedBy)
	}
	if len(r.Edges.Menus) > 0 {
		menuIds := make([]int64, len(r.Edges.Menus))
		for i, m := range r.Edges.Menus {
			menuIds[i] = m.ID
		}
		resp.MenuIds = menuIds
	}
	return resp
}

func roleApiIds(enf *casbinapi.SyncedEnforcer, roleCode string, tenantId int64) []int64 {
	dom := strconv.FormatInt(tenantId, 10)
	policies, _ := enf.GetFilteredPolicy(0, roleCode, dom)
	ids := make([]int64, 0, len(policies))
	for _, p := range policies {
		if len(p) > 4 {
			if id, err := strconv.ParseInt(p[4], 10, 64); err == nil {
				ids = append(ids, id)
			}
		}
	}
	return ids
}
