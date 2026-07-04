package syslogslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysoperationlog"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperationLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperationLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperationLogListLogic {
	return &GetOperationLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOperationLogListLogic) GetOperationLogList(in *apps.LogPageReq) (*apps.OperationLogListResp, error) {
	query := l.svcCtx.DB.SysOperationLog.Query()

	if in.Module != nil && *in.Module != "" {
		query.Where(sysoperationlog.ModuleContains(*in.Module))
	}
	if in.Operation != nil && *in.Operation != "" {
		query.Where(sysoperationlog.OperationContains(*in.Operation))
	}
	if in.OperatorName != nil && *in.OperatorName != "" {
		query.Where(sysoperationlog.OperatorNameContains(*in.OperatorName))
	}
	if in.Path != nil && *in.Path != "" {
		query.Where(sysoperationlog.PathContains(*in.Path))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.OperationLogListResp{Code: 500, Msg: fmt.Sprintf("查询操作日志总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	logs, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Desc(sysoperationlog.FieldID)).
		All(l.ctx)
	if err != nil {
		return &apps.OperationLogListResp{Code: 500, Msg: fmt.Sprintf("查询操作日志列表失败: %v", err)}, nil
	}

	list := make([]*apps.OperationLog, 0, len(logs))
	for _, l := range logs {
		list = append(list, operationLogToPb(l))
	}

	return &apps.OperationLogListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}

func operationLogToPb(l *ent.SysOperationLog) *apps.OperationLog {
	log := &apps.OperationLog{
		Id:            &l.ID,
		IdStr:         strPtr(fmt.Sprintf("%d", l.ID)),
		Module:        &l.Module,
		Operation:     &l.Operation,
		RequestMethod: &l.Method,
		RequestUrl:    &l.Path,
		RequestParam:  &l.Params,
		ResponseData:  &l.Result,
		Duration:      &l.Duration,
		OperatorId:    &l.OperatorID,
		OperatorIdStr: strPtr(fmt.Sprintf("%d", l.OperatorID)),
		OperatorName:  &l.OperatorName,
		OperatorIp:    &l.IP,
		TenantId:      &l.TenantID,
		TenantIdStr:   strPtr(fmt.Sprintf("%d", l.TenantID)),
	}

	return log
}