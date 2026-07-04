package sysdict_dataslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdict"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByDictKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictDataByDictKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByDictKeyLogic {
	return &GetDictDataByDictKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictDataByDictKeyLogic) GetDictDataByDictKey(in *apps.DictReq) (*apps.DictDataListResp, error) {
	if in.Key == nil || *in.Key == "" {
		return &apps.DictDataListResp{Code: 500, Msg: "字典键不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)

	dictDatas, err := l.svcCtx.DB.SysDictData.TenantAwareQuery(tenantId).
		Where(sysdictdata.HasSysDictWith(sysdict.KeyEQ(*in.Key))).
		All(l.ctx)
	if err != nil {
		return &apps.DictDataListResp{Code: 500, Msg: fmt.Sprintf("获取字典数据失败: %v", err)}, nil
	}

	list := make([]*apps.DictData, 0, len(dictDatas))
	for _, dd := range dictDatas {
		list = append(list, dictDataToPb(dd))
	}

	return &apps.DictDataListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(len(list)),
	}, nil
}