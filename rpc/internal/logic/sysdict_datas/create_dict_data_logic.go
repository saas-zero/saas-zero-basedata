package sysdict_dataslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysdictdata"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateDictDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDictDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictDataLogic {
	return &CreateDictDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDictDataLogic) CreateDictData(in *apps.DictDataReq) (*apps.DictDataResp, error) {
	if in.DictId == nil || *in.DictId <= 0 {
		return &apps.DictDataResp{Code: 500, Msg: "字典ID不能为空"}, nil
	}

	if in.Name == nil || *in.Name == "" {
		return &apps.DictDataResp{Code: 500, Msg: "字典数据名称不能为空"}, nil
	}

	if in.Key == nil || *in.Key == "" {
		return &apps.DictDataResp{Code: 500, Msg: "字典数据键不能为空"}, nil
	}

	if in.Value == nil || *in.Value == "" {
		return &apps.DictDataResp{Code: 500, Msg: "字典数据值不能为空"}, nil
	}

	create := l.svcCtx.DB.SysDictData.Create().
		SetDictID(*in.DictId).
		SetName(*in.Name).
		SetKey(*in.Key).
		SetValue(*in.Value)

	if in.Status != nil {
		create.SetStatus(sysdictdata.Status(*in.Status))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}

	dictData, err := create.Save(l.ctx)
	if err != nil {
		return &apps.DictDataResp{Code: 500, Msg: fmt.Sprintf("创建字典数据失败: %v", err)}, nil
	}

	return &apps.DictDataResp{
		Code: 200,
		Msg:  "创建成功",
		Data: dictDataToPb(dictData),
	}, nil
}