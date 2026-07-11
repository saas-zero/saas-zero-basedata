package sysuserslogic

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *apps.UserReq) (*apps.EmptyResp, error) {
	userId := mixins.GetCurrentUserId(l.ctx)
	userName := mixins.GetCurrentUserName(l.ctx)
	ctx := mixins.SetCurrentUserId(l.ctx, userId)
	ctx = mixins.SetCurrentUserName(ctx, userName)

	hash, err := bcrypt.GenerateFromPassword([]byte(in.GetPassword()), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	err = l.svcCtx.DB.SysUser.UpdateOneID(in.GetId()).
		SetPassword(string(hash)).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &apps.EmptyResp{Code: int32(errno.Success.Code), Msg: errno.Success.Msg}, nil
}
