package syspackageslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/syspackage"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePackageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePackageLogic {
	return &UpdatePackageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePackageLogic) UpdatePackage(in *apps.PackageReq) (*apps.PackageResp, error) {
	if in.Id == nil || *in.Id <= 0 {
		return &apps.PackageResp{Code: 500, Msg: "套餐ID不能为空"}, nil
	}

	update := l.svcCtx.DB.SysPackage.UpdateOneID(int64(*in.Id))

	if in.Name != nil {
		update.SetName(*in.Name)
	}
	if in.Code != nil {
		update.SetCode(*in.Code)
	}
	if in.Status != nil {
		update.SetStatus(syspackage.Status(*in.Status))
	}
	if in.Sort != nil {
		update.SetSort(uint32(*in.Sort))
	}
	if in.Remark != nil {
		update.SetRemark(*in.Remark)
	}

	pkg, err := update.Save(l.ctx)
	if err != nil {
		return &apps.PackageResp{Code: 500, Msg: fmt.Sprintf("更新套餐失败: %v", err)}, nil
	}

	return &apps.PackageResp{
		Code: 200,
		Msg:  "更新成功",
		Data: packageToPb(pkg),
	}, nil
}