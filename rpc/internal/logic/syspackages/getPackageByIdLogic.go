package syspackageslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPackageByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPackageByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPackageByIdLogic {
	return &GetPackageByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPackageByIdLogic) GetPackageById(in *apps.IdReq) (*apps.PackageResp, error) {
	// todo: add your logic here and delete this line

	return &apps.PackageResp{}, nil
}
