package sysapislogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysapi"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApiLogic) UpdateApi(in *apps.ApiReq) (*apps.ApiResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.ApiResp{Code: 500, Msg: "API ID不能为空"}, nil
	}

	update := l.svcCtx.DB.SysApi.UpdateOneID(int64(*in.Id))

	if in.ApiName != nil {
		update.SetAPIName(*in.ApiName)
	}
	if in.ApiType != nil {
		update.SetAPIType(sysapi.APIType(*in.ApiType))
	}
	if in.ApiPath != nil {
		update.SetAPIPath(*in.ApiPath)
	}
	if in.ApiMethod != nil {
		update.SetAPIMethod(sysapi.APIMethod(*in.ApiMethod))
	}
	if in.Status != nil {
		update.SetStatus(sysapi.Status(*in.Status))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	api, err := update.Save(l.ctx)
	if err != nil {
		return &apps.ApiResp{Code: 500, Msg: fmt.Sprintf("更新API失败: %v", err)}, nil
	}

	return &apps.ApiResp{
		Code: 200,
		Msg:  "更新成功",
		Data: apiToPb(api),
	}, nil
}