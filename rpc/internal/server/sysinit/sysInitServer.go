package server

import (
	"context"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/logic/sysinit"
	"github.com/saas-zero/saas-zero-basedata/rpc/internal/svc"
)

type SysInitServer struct {
	svcCtx *svc.ServiceContext
	apps.UnimplementedSysInitServer
}

func NewSysInitServer(svcCtx *svc.ServiceContext) *SysInitServer {
	return &SysInitServer{
		svcCtx: svcCtx,
	}
}

func (s *SysInitServer) InitAll(ctx context.Context, in *apps.EmptyReq) (*apps.EmptyResp, error) {
	l := sysinitlogic.NewInitAllLogic(ctx, s.svcCtx)
	return l.InitAll(in)
}
