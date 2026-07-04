package sysdictslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdict"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

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

	if in.Name != nil && *in.Name != "" {
		query.Where(sysdict.NameContains(*in.Name))
	}
	if in.Key != nil && *in.Key != "" {
		query.Where(sysdict.KeyContains(*in.Key))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysdict.StatusEQ(sysdict.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.DictListResp{Code: 500, Msg: fmt.Sprintf("查询字典总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	dicts, err := query.
		Offset((page - 1) * size).
		Limit(size).
		All(l.ctx)
	if err != nil {
		return &apps.DictListResp{Code: 500, Msg: fmt.Sprintf("查询字典列表失败: %v", err)}, nil
	}

	list := make([]*apps.Dict, 0, len(dicts))
	for _, d := range dicts {
		list = append(list, dictToPb(d, tenantId))
	}

	return &apps.DictListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}