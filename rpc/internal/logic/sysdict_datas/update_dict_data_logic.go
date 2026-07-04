package sysdict_dataslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictDataLogic {
	return &UpdateDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDictDataLogic) UpdateDictData(in *apps.DictDataReq) (*apps.DictDataResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.DictDataResp{Code: 500, Msg: "字典数据ID不能为空"}, nil
	}

	update := l.svcCtx.DB.SysDictData.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.Key != nil {
		update.SetKey(*in.Key)
	}
	if in.Value != nil {
		update.SetValue(*in.Value)
	}
	if in.Status != nil {
		update.SetStatus(sysdictdata.Status(*in.Status))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	dictData, err := update.Save(l.ctx)
	if err != nil {
		return &apps.DictDataResp{Code: 500, Msg: fmt.Sprintf("更新字典数据失败: %v", err)}, nil
	}

	return &apps.DictDataResp{
		Code: 200,
		Msg:  "更新成功",
		Data: dictDataToPb(dictData),
	}, nil
}