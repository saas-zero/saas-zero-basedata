package systenantslogic

import (
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"google.golang.org/protobuf/proto"
)

func tenantToResp(t *ent.SysTenant) *apps.Tenant {
	resp := &apps.Tenant{
		Id:        proto.Int64(t.ID),
		IdStr:     proto.String(strconv.FormatInt(t.ID, 10)),
		Name:      proto.String(t.Name),
		Code:      proto.String(t.Code),
		Status:    proto.String(string(t.Status)),
		Remark:    proto.String(t.Remark),
		CreatedAt: proto.Int64(t.CreatedAt.UnixMilli()),
		UpdatedAt: proto.Int64(t.UpdatedAt.UnixMilli()),
	}
	if t.CreatedBy != "" {
		resp.CreatedBy = proto.String(t.CreatedBy)
	}
	if t.UpdatedBy != "" {
		resp.UpdatedBy = proto.String(t.UpdatedBy)
	}
	if t.AdminID > 0 {
		resp.AdminId = proto.Int64(t.AdminID)
		resp.AdminIdStr = proto.String(strconv.FormatInt(t.AdminID, 10))
	}
	if t.ParentID > 0 {
		resp.ParentId = proto.Int64(t.ParentID)
		resp.ParentIdStr = proto.String(strconv.FormatInt(t.ParentID, 10))
	}
	if t.PackageID > 0 {
		resp.PackageId = proto.Int64(t.PackageID)
		resp.PackageIdStr = proto.String(strconv.FormatInt(t.PackageID, 10))
	}
	if t.Edges.SysPackage != nil {
		resp.PackageName = proto.String(t.Edges.SysPackage.Name)
	}
	if !t.ExpiredAt.IsZero() {
		resp.ExpiredAt = proto.Int64(t.ExpiredAt.UnixMilli())
	}
	return resp
}
