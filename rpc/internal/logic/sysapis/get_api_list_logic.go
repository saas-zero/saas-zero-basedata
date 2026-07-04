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

type GetApiListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetApiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApiListLogic {
	return &GetApiListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetApiListLogic) GetApiList(in *apps.ApiPageReq) (*apps.ApiListResp, error) {
	query := l.svcCtx.DB.SysApi.Query()

	if in.ApiName != nil && *in.ApiName != "" {
		query.Where(sysapi.APINameContains(*in.ApiName))
	}
	if in.ApiPath != nil && *in.ApiPath != "" {
		query.Where(sysapi.APIPathContains(*in.ApiPath))
	}
	if in.ApiType != nil && *in.ApiType != "" {
		query.Where(sysapi.APITypeEQ(sysapi.APIType(*in.ApiType)))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysapi.StatusEQ(sysapi.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.ApiListResp{Code: 500, Msg: fmt.Sprintf("查询API总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	apis, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Desc(sysapi.FieldCreatedAt)).
		All(l.ctx)
	if err != nil {
		return &apps.ApiListResp{Code: 500, Msg: fmt.Sprintf("查询API列表失败: %v", err)}, nil
	}

	list := make([]*apps.Api, 0, len(apis))
	for _, a := range apis {
		list = append(list, apiToPb(a))
	}

	return &apps.ApiListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}