package syspackageslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/syspackage"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePackageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePackageLogic {
	return &CreatePackageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePackageLogic) CreatePackage(in *apps.PackageReq) (*apps.PackageResp, error) {
	if in.Name == nil || *in.Name == "" {
		return &apps.PackageResp{Code: 500, Msg: "套餐名称不能为空"}, nil
	}

	if in.Code == nil || *in.Code == "" {
		return &apps.PackageResp{Code: 500, Msg: "套餐编码不能为空"}, nil
	}

	create := l.svcCtx.DB.SysPackage.Create().
		SetName(*in.Name).
		SetCode(*in.Code)

	if in.Status != nil {
		create.SetStatus(syspackage.Status(*in.Status))
	}
	if in.Sort != nil {
		create.SetSort(uint32(*in.Sort))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}

	pkg, err := create.Save(l.ctx)
	if err != nil {
		return &apps.PackageResp{Code: 500, Msg: fmt.Sprintf("创建套餐失败: %v", err)}, nil
	}

	return &apps.PackageResp{
		Code: 200,
		Msg:  "创建成功",
		Data: packageToPb(pkg),
	}, nil
}