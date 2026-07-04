package syspackageslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPackageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPackageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPackageListLogic {
	return &GetPackageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPackageListLogic) GetPackageList(in *apps.PackagePageReq) (*apps.PackageListResp, error) {
	// todo: add your logic here and delete this line

	return &apps.PackageListResp{}, nil
}
