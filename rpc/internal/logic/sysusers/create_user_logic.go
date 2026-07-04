package sysuserslogic

import (
	"context"
	"fmt"

	"github.com/saas-zero/saas-zero-basedata/ent/sysuser"
	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateUserLogic) CreateUser(in *apps.UserReq) (*apps.UserResp, error) {
	if in.Username == nil || *in.Username == "" {
		return &apps.UserResp{Code: 500, Msg: "用户名不能为空"}, nil
	}

	password := "123456"
	if in.Password != nil && *in.Password != "" {
		password = *in.Password
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &apps.UserResp{Code: 500, Msg: fmt.Sprintf("密码加密失败: %v", err)}, nil
	}

	create := l.svcCtx.DB.SysUser.Create().
		SetUsername(*in.Username).
		SetPassword(string(hashedPassword))

	if in.Nickname != nil {
		create.SetNickname(*in.Nickname)
	}
	if in.Mobile != nil {
		create.SetMobile(*in.Mobile)
	}
	if in.Email != nil {
		create.SetEmail(*in.Email)
	}
	if in.DeptId != nil {
		create.SetDeptID(*in.DeptId)
	}
	if in.Status != nil {
		create.SetStatus(sysuser.Status(*in.Status))
	}
	if in.Remark != nil {
		create.SetRemark(*in.Remark)
	}
	if len(in.RoleIds) > 0 {
		create.AddRoleIDs(in.RoleIds...)
	}

	user, err := create.Save(l.ctx)
	if err != nil {
		return &apps.UserResp{Code: 500, Msg: fmt.Sprintf("创建用户失败: %v", err)}, nil
	}

	return &apps.UserResp{
		Code: 200,
		Msg:  "创建成功",
		Data: userToPb(user),
	}, nil
}