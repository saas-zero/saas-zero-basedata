package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
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
	if in.Id == nil || *in.Id <= 0 {
		return &apps.EmptyResp{Code: 500, Msg: "用户ID不能为空"}, nil
	}

	password := "123456"
	if in.Password != nil && *in.Password != "" {
		password = *in.Password
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("密码加密失败: %v", err)}, nil
	}

	_, err = l.svcCtx.DB.SysUser.UpdateOneID(int64(*in.Id)).
		SetPassword(string(hashedPassword)).
		Save(l.ctx)
	if err != nil {
		return &apps.EmptyResp{Code: 500, Msg: fmt.Sprintf("重置密码失败: %v", err)}, nil
	}

	return &apps.EmptyResp{Code: 200, Msg: "密码重置成功"}, nil
}