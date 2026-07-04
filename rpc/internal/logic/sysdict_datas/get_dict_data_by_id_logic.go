package sysdict_dataslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDictDataByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDictDataByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDictDataByIdLogic {
	return &GetDictDataByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDictDataByIdLogic) GetDictDataById(in *apps.IdReq) (*apps.DictDataResp, error) {
	dd, err := l.svcCtx.DB.SysDictData.Query().
		Where(sysdictdata.IDEQ(int64(in.Id))).
		Only(l.ctx)
	if err != nil {
		return &apps.DictDataResp{Code: 500, Msg: fmt.Sprintf("获取字典数据失败: %v", err)}, nil
	}

	return &apps.DictDataResp{
		Code: 200,
		Msg:  "success",
		Data: dictDataToPb(dd),
	}, nil
}

func dictDataToPb(dd *ent.SysDictData) *apps.DictData {
	dictData := &apps.DictData{
		Id:       &dd.ID,
		IdStr:    strPtr(fmt.Sprintf("%d", dd.ID)),
		DictId:   &dd.DictID,
		DictIdStr: strPtr(fmt.Sprintf("%d", dd.DictID)),
		Name:     &dd.Name,
		Key:      &dd.Key,
		Value:    &dd.Value,
		Status:   strPtr(string(dd.Status)),
		Remark:   &dd.Remark,
		TenantId: &dd.TenantID,
		TenantIdStr: strPtr(fmt.Sprintf("%d", dd.TenantID)),
	}

	createdAt := dd.CreatedAt.Unix()
	dictData.CreatedAt = &createdAt
	dictData.CreatedBy = &dd.CreatedBy

	updatedAt := dd.UpdatedAt.Unix()
	dictData.UpdatedAt = &updatedAt
	dictData.UpdatedBy = &dd.UpdatedBy

	return dictData
}

func strPtr(s string) *string {
	return &s
}