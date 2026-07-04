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
	pkg, err := l.svcCtx.DB.SysPackage.Query().
		Where(syspackage.IDEQ(int64(in.Id))).
		Only(l.ctx)
	if err != nil {
		return &apps.PackageResp{Code: 500, Msg: fmt.Sprintf("获取套餐失败: %v", err)}, nil
	}

	return &apps.PackageResp{
		Code: 200,
		Msg:  "success",
		Data: packageToPb(pkg),
	}, nil
}

func packageToPb(p *ent.SysPackage) *apps.Package {
	sort := int32(p.Sort)
	pkg := &apps.Package{
		Id:       &p.ID,
		IdStr:    strPtr(fmt.Sprintf("%d", p.ID)),
		Name:     &p.Name,
		Code:     &p.Code,
		Status:   strPtr(string(p.Status)),
		Sort:     &sort,
		Remark:   &p.Remark,
	}

	createdAt := p.CreatedAt.Unix()
	pkg.CreatedAt = &createdAt
	pkg.CreatedBy = &p.CreatedBy

	updatedAt := p.UpdatedAt.Unix()
	pkg.UpdatedAt = &updatedAt
	pkg.UpdatedBy = &p.UpdatedBy

	return pkg
}

func strPtr(s string) *string {
	return &s
}