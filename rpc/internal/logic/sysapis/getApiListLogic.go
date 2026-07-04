package sysapislogic

import (
	"context"

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
	query := l.svcCtx.DB.SysApi.ActiveQuery()
	if in.GetApiName() != "" {
		query = query.Where(sysapi.APINameContains(in.GetApiName()))
	}
	if in.GetApiPath() != "" {
		query = query.Where(sysapi.APIPathContains(in.GetApiPath()))
	}
	if in.GetApiType() != "" {
		query = query.Where(sysapi.APITypeEQ(sysapi.APIType(in.GetApiType())))
	}
	if in.GetStatus() != "" {
		query = query.Where(sysapi.StatusEQ(sysapi.Status(in.GetStatus())))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return nil, err
	}

	page := int(in.GetPage())
	size := int(in.GetSize())
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	apis, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysapi.FieldCreatedAt)).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.Api, len(apis))
	for i, a := range apis {
		list[i] = apiToResp(a)
	}
	return &apps.ApiListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}
