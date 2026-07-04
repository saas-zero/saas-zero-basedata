package syspackageslogic

import (
	"context"
	"fmt"
	"time"

	"github.com/saas-zero/saas-zero-basedata/ent/systenant"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePackageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePackageLogic {
	return &DeletePackageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePackageLogic) DeletePackage(in *apps.IdsReq) (*apps.EmptyResp, error) {
	if len(in.Ids) == 0 {
		return &apps.EmptyResp{Code: 500, Msg: "请选择要删除的套餐"}, nil
	}

	for _, id := range in.Ids {
		count, err := l.svcCtx.DB.SysTenant.Query().
			Where(systenant.PackageIDEQ(id)).
			Count(l.ctx)
		if err != nil {
			return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("查询租户关联失败: %v", err)}, nil
		}

		if count > 0 {
			return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("该套餐已被 %d 个租户使用，无法删除", count)}, nil
		}

		_, err = l.svcCtx.DB.SysPackage.UpdateOneID(id).
			SetDeletedAt(time.Now()).
			Save(l.ctx)
		if err != nil {
			return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("删除套餐失败: %v", err)}, nil
		}
	}

	return &apps.EmptyResp{Code: 200, Msg: "删除成功"}, nil
}