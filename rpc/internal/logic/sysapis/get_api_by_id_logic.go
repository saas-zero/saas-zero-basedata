package sysapislogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysapi"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApiByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiByIdLogic {
	return &GetApiByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApiByIdLogic) GetApiById(in *apps.IdReq) (*apps.ApiResp, error) {
	api, err := l.svcCtx.DB.SysApi.Query().
		Where(sysapi.IDEQ(int64(in.Id))).
		Only(l.ctx)
	if err != nil {
		return &apps.ApiResp{Code: 500, Msg: fmt.Sprintf("获取API失败: %v", err)}, nil
	}

	return &apps.ApiResp{
		Code: 200,
		Msg:  "success",
		Data: apiToPb(api),
	}, nil
}

func apiToPb(a *ent.SysApi) *apps.Api {
	api := &apps.Api{
		Id:       &a.ID,
		IdStr:    strPtr(fmt.Sprintf("%d", a.ID)),
		ApiName:  &a.APIName,
		ApiType:  strPtr(string(a.APIType)),
		ApiPath:  &a.APIPath,
		Status:   strPtr(string(a.Status)),
		Remark:   &a.Remark,
	}

	if a.APIMethod != "" {
		api.ApiMethod = strPtr(string(a.APIMethod))
	}

	createdAt := a.CreatedAt.Unix()
	api.CreatedAt = &createdAt
	api.CreatedBy = &a.CreatedBy

	updatedAt := a.UpdatedAt.Unix()
	api.UpdatedAt = &updatedAt
	api.UpdatedBy = &a.UpdatedBy

	return api
}

func strPtr(s string) *string {
	return &s
}