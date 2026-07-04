package sysapislogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysapi"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateApiLogic) CreateApi(in *apps.ApiReq) (*apps.ApiResp, error) {
	if in.ApiName == nil || *in.ApiName == "" {
		return &apps.ApiResp{Code: 500, Msg: "API名称不能为空"}, nil
	}

	if in.ApiPath == nil || *in.ApiPath == "" {
		return &apps.ApiResp{Code: 500, Msg: "API路径不能为空"}, nil
	}

	create := l.svcCtx.DB.SysApi.Create().
		SetAPIName(*in.ApiName).
		SetAPIPath(*in.ApiPath)

	if in.ApiType != nil {
		create.SetAPIType(sysapi.APIType(*in.ApiType))
	}
	if in.ApiMethod != nil {
		create.SetAPIMethod(sysapi.APIMethod(*in.ApiMethod))
	}
	if in.Status != nil {
		create.SetStatus(sysapi.Status(*in.Status))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}

	api, err := create.Save(l.ctx)
	if err != nil {
		return &apps.ApiResp{Code: 500, Msg: fmt.Sprintf("创建API失败: %v", err)}, nil
	}

	return &apps.ApiResp{
		Code: 200,
		Msg:  "创建成功",
		Data: apiToPb(api),
	}, nil
}