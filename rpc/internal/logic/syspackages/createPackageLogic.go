package syspackageslogic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &apps.PackageResp{}, nil
}
