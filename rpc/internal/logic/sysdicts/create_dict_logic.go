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

type CreateDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateDictLogic {
	return &CreateDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateDictLogic) CreateDict(in *apps.DictReq) (*apps.DictResp, error) {
	if in.Name == nil || *in.Name == "" {
		return &apps.DictResp{Code: 500, Msg: "字典名称不能为空"}, nil
	}

	if in.Key == nil || *in.Key == "" {
		return &apps.DictResp{Code: 500, Msg: "字典键不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)

	create := l.svcCtx.DB.SysDict.Create().
		SetName(*in.Name).
		SetKey(*in.Key)

	if in.Status != nil {
		create.SetStatus(sysdict.Status(*in.Status))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}

	dict, err := create.Save(l.ctx)
	if err != nil {
		return &apps.DictResp{Code: 500, Msg: fmt.Sprintf("创建字典失败: %v", err)}, nil
	}

	return &apps.DictResp{
		Code: 200,
		Msg:  "创建成功",
		Data: dictToPb(dict, tenantId),
	}, nil
}