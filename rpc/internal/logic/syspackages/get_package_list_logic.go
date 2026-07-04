package syspackageslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent"
	"github.com/saas-zero/saas-zero-basedata/ent/syspackage"
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
	query := l.svcCtx.DB.SysPackage.Query()

	if in.Name != nil && *in.Name != "" {
		query.Where(syspackage.NameContains(*in.Name))
	}
	if in.Code != nil && *in.Code != "" {
		query.Where(syspackage.CodeContains(*in.Code))
	}
	if in.Status != nil && *in.Status != "" {
		query.Where(syspackage.StatusEQ(syspackage.Status(*in.Status)))
	}

	total, err := query.Count(l.ctx)
	if err != nil {
		return &apps.PackageListResp{Code: 500, Msg: fmt.Sprintf("查询套餐总数失败: %v", err)}, nil
	}

	page := int(in.Page)
	size := int(in.Size)
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}

	pkgs, err := query.
		Offset((page - 1) * size).
		Limit(size).
		Order(ent.Asc(syspackage.FieldSort)).
		All(l.ctx)
	if err != nil {
		return &apps.PackageListResp{Code: 500, Msg: fmt.Sprintf("查询套餐列表失败: %v", err)}, nil
	}

	list := make([]*apps.Package, 0, len(pkgs))
	for _, p := range pkgs {
		list = append(list, packageToPb(p))
	}

	return &apps.PackageListResp{
		Code:  200,
		Msg:   "success",
		List:  list,
		Total: int64(total),
	}, nil
}