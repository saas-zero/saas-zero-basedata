package sysdictslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdict"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictListLogic {
	return &GetDictListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictListLogic) GetDictList(in *apps.DictPageReq) (*apps.DictListResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	query := l.svcCtx.DB.SysDict.TenantAwareQuery(tenantId)
	if in.GetName() != "" {
		query = query.Where(sysdict.NameContains(in.GetName()))
	}
	if in.GetKey() != "" {
		query = query.Where(sysdict.KeyContains(in.GetKey()))
	}
	if in.GetStatus() != "" {
		query = query.Where(sysdict.StatusEQ(sysdict.Status(in.GetStatus())))
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

	dicts, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysdict.FieldCreatedAt)).
		WithSysDictDatas().
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.Dict, len(dicts))
	for i, d := range dicts {
		list[i] = dictToResp(d)
	}
	return &apps.DictListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}
