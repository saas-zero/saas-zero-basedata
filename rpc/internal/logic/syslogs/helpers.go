package syslogslogic

import (
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"google.golang.org/protobuf/proto"
)

func loginLogToResp(l *ent.SysLoginLog) *apps.LoginLog {
	return &apps.LoginLog{
		Id:       proto.Int64(l.ID),
		IdStr:    proto.String(strconv.FormatInt(l.ID, 10)),
		Username: proto.String(l.Username),
		LoginIp:  proto.String(l.IP),
		Status:   proto.String(string(l.Status)),
		Msg:      proto.String(l.Message),
		LoginAt:  proto.Int64(l.LoginTime.UnixMilli()),
	}
}

func operationLogToResp(o *ent.SysOperationLog) *apps.OperationLog {
	return &apps.OperationLog{
		Id:            proto.Int64(o.ID),
		IdStr:         proto.String(strconv.FormatInt(o.ID, 10)),
		Module:        proto.String(o.Module),
		Operation:     proto.String(o.Operation),
		RequestMethod: proto.String(o.Method),
		RequestUrl:    proto.String(o.Path),
		RequestParam:  proto.String(o.Params),
		ResponseData:  proto.String(o.Result),
		Duration:      proto.Int64(o.Duration),
		Status:        proto.String("success"),
		OperatorId:    proto.Int64(o.OperatorID),
		OperatorIdStr: proto.String(strconv.FormatInt(o.OperatorID, 10)),
		OperatorName:  proto.String(o.OperatorName),
		OperatorIp:    proto.String(o.IP),
		TenantId:      proto.Int64(o.TenantID),
		TenantIdStr:   proto.String(strconv.FormatInt(o.TenantID, 10)),
	}
}
