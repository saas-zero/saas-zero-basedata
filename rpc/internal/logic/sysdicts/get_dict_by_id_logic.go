package sysdictslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdict"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictByIdLogic {
	return &GetDictByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictByIdLogic) GetDictById(in *apps.IdReq) (*apps.DictResp, error) {
	tenantId := mixins.GetCurrentTenantId(l.ctx)

	dict, err := l.svcCtx.DB.SysDict.Query().
		Where(sysdict.IDEQ(int64(in.Id))).
		WithSysDictDatas().
		Only(l.ctx)
	if err != nil {
		return &apps.DictResp{Code: 500, Msg: fmt.Sprintf("获取字典失败: %v", err)}, nil
	}

	return &apps.DictResp{
		Code: 200,
		Msg:  "success",
		Data: dictToPb(dict, tenantId),
	}, nil
}

func dictToPb(d *ent.SysDict, tenantId int64) *apps.Dict {
	dict := &apps.Dict{
		Id:         &d.ID,
		IdStr:      strPtr(fmt.Sprintf("%d", d.ID)),
		Name:       &d.Name,
		Key:        &d.Key,
		Status:     strPtr(string(d.Status)),
		Remark:     &d.Remark,
		TenantId:   &d.TenantID,
		TenantIdStr: strPtr(fmt.Sprintf("%d", d.TenantID)),
	}

	if len(d.Edges.SysDictDatas) > 0 {
		dataList := make([]*apps.DictData, 0, len(d.Edges.SysDictDatas))
		for _, dd := range d.Edges.SysDictDatas {
			dataList = append(dataList, dictDataToPb(dd))
		}
		dict.DictData = dataList
	}

	createdAt := d.CreatedAt.Unix()
	dict.CreatedAt = &createdAt
	dict.CreatedBy = &d.CreatedBy

	updatedAt := d.UpdatedAt.Unix()
	dict.UpdatedAt = &updatedAt
	dict.UpdatedBy = &d.UpdatedBy

	return dict
}

func dictDataToPb(dd *ent.SysDictData) *apps.DictData {
	data := &apps.DictData{
		Id:         &dd.ID,
		IdStr:      strPtr(fmt.Sprintf("%d", dd.ID)),
		DictId:     &dd.DictID,
		DictIdStr:  strPtr(fmt.Sprintf("%d", dd.DictID)),
		Name:       &dd.Name,
		Key:        &dd.Key,
		Value:      &dd.Value,
		Status:     strPtr(string(dd.Status)),
		Remark:     &dd.Remark,
		TenantId:   &dd.TenantID,
		TenantIdStr: strPtr(fmt.Sprintf("%d", dd.TenantID)),
	}

	createdAt := dd.CreatedAt.Unix()
	data.CreatedAt = &createdAt
	data.CreatedBy = &dd.CreatedBy

	updatedAt := dd.UpdatedAt.Unix()
	data.UpdatedAt = &updatedAt
	data.UpdatedBy = &dd.UpdatedBy

	return data
}

func strPtr(s string) *string {
	return &s
}