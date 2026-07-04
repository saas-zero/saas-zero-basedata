package sysdictdataslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictDataListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataListLogic {
	return &GetDictDataListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictDataListLogic) GetDictDataList(in *apps.DictDataPageReq) (*apps.DictDataListResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	query := l.svcCtx.DB.SysDictData.TenantAwareQuery(tenantId)

	if in.GetDictId() > 0 {
		query = query.Where(sysdictdata.DictIDEQ(in.GetDictId()))
	}
	if in.GetKey() != "" {
		query = query.Where(sysdictdata.KeyContains(in.GetKey()))
	}
	if in.GetValue() != "" {
		query = query.Where(sysdictdata.ValueContains(in.GetValue()))
	}
	if in.GetStatus() != "" {
		query = query.Where(sysdictdata.StatusEQ(sysdictdata.Status(in.GetStatus())))
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

	items, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysdictdata.FieldCreatedAt)).
		All(l.ctx)
	if err != nil {
		return nil, err
	}

	list := make([]*apps.DictData, len(items))
	for i, d := range items {
		list[i] = dictDataToResp(d)
	}
	return &apps.DictDataListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}
