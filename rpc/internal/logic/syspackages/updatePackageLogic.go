package syspackageslogic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &apps.PackageResp{}, nil
}
