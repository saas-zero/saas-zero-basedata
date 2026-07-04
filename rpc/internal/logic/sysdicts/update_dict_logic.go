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

type UpdateDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateDictLogic {
	return &UpdateDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateDictLogic) UpdateDict(in *apps.DictReq) (*apps.DictResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.DictResp{Code: 500, Msg: "字典ID不能为空"}, nil
	}

	tenantId := mixins.GetCurrentTenantId(l.ctx)

	update := l.svcCtx.DB.SysDict.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.Key != nil {
		update.SetKey(*in.Key)
	}
	if in.Status != nil {
		update.SetStatus(sysdict.Status(*in.Status))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	dict, err := update.Save(l.ctx)
	if err != nil {
		return &apps.DictResp{Code: 500, Msg: fmt.Sprintf("更新字典失败: %v", err)}, nil
	}

	return &apps.DictResp{
		Code: 200,
		Msg:  "更新成功",
		Data: dictToPb(dict, tenantId),
	}, nil
}