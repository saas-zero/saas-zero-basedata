package syspackageslogic

import (
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"google.golang.org/protobuf/proto"
)

func packageToResp(p *ent.SysPackage) *apps.Package {
	resp := &apps.Package{
		Id:        proto.Int64(p.ID),
		IdStr:     proto.String(strconv.FormatInt(p.ID, 10)),
		Name:      proto.String(p.Name),
		Code:      proto.String(p.Code),
		Status:    proto.String(string(p.Status)),
		Sort:      proto.Int32(int32(p.Sort)),
		Remark:    proto.String(p.Remark),
		CreatedAt: proto.Int64(p.CreatedAt.UnixMilli()),
		UpdatedAt: proto.Int64(p.UpdatedAt.UnixMilli()),
	}
	if p.CreatedBy != "" {
		resp.CreatedBy = proto.String(p.CreatedBy)
	}
	if p.UpdatedBy != "" {
		resp.UpdatedBy = proto.String(p.UpdatedBy)
	}
	return resp
}
