package sysdict_dataslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

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

	if in.DictId != nil && *in.DictId > 0 {
		query.Where(sysdictdata.DictIDEQ(*in.DictId))
	}
	if in.Key != nil && *in.Key != "" {
		query.Where(sysdictdata.KeyContains(*in.Key))
	}
	if in.Value != nil && *in.Value != "" {
		query.Where(sysdictdata.ValueContains(*in.Value))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(sysdictdata.StatusEQ(sysdictdata.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.DictDataListResp{Code: 500, Msg: fmt.Sprintf("查询字典数据总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	dictDatas, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(sysdictdata.FieldCreatedAt)).
		All(l.ctx)
	if err != nil {
		return &apps.DictDataListResp{Code: 500, Msg: fmt.Sprintf("查询字典数据列表失败: %v", err)}, nil
	}

	list := make([]*apps.DictData, 0, len(dictDatas))
	for _, dd := range dictDatas {
		list = append(list, dictDataToPb(dd))
	}

	return &apps.DictDataListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}