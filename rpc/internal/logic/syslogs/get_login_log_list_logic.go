package syslogslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysloginlog"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLoginLogListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLoginLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLoginLogListLogic {
	return &GetLoginLogListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLoginLogListLogic) GetLoginLogList(in *apps.LogPageReq) (*apps.LoginLogListResp, error) {
	query := l.svcCtx.DB.SysLoginLog.Query()

	if in.Username != nil && *in.Username != "" {
		query.Where(sysloginlog.UsernameContains(*in.Username))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysloginlog.StatusEQ(sysloginlog.Status(*in.Status)))
	}
	if in.Ip != nil && *in.Ip != "" {
		query.Where(sysloginlog.IPContains(*in.Ip))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.LoginLogListResp{Code: 500, Msg: fmt.Sprintf("查询登录日志总数失败: %v", err)}, nil
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
		Order(ent.Desc(sysloginlog.FieldLoginTime)).
		All(l.ctx)
	if err != nil {
		return &apps.LoginLogListResp{Code: 500, Msg: fmt.Sprintf("查询登录日志列表失败: %v", err)}, nil
	}

	list := make([]*apps.LoginLog, 0, len(logs))
	for _, l := range logs {
		list = append(list, loginLogToPb(l))
	}

	return &apps.LoginLogListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}

func loginLogToPb(l *ent.SysLoginLog) *apps.LoginLog {
	loginAt := l.LoginTime.Unix()
	log := &apps.LoginLog{
		Id:       &l.ID,
		IdStr:    strPtr(fmt.Sprintf("%d", l.ID)),
		Username: &l.Username,
		LoginIp:  &l.IP,
		Status:   strPtr(string(l.Status)),
		Msg:      &l.Message,
		LoginAt:  &loginAt,
	}

	return log
}

func strPtr(s string) *string {
	return &s
}