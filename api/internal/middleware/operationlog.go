package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/protobuf/proto"
)

// OperationLog 记录所有写操作（非 GET）的审计日志。跳过 GET 与 /init/* 路由。
// 操作人/租户信息来自 JwtAuth 已注入的 context，因此必须注册在 JwtAuth 之后。
func OperationLog(cli apps.SysLogsClient) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || strings.HasPrefix(r.URL.Path, "/init/") {
				next(w, r)
				return
			}

			start := time.Now()
			next(w, r)
			duration := time.Since(start).Milliseconds()

			ctx := r.Context()
			module, operation := parseModuleOperation(r.URL.Path)

			log := &apps.OperationLog{
				Module:        proto.String(module),
				Operation:     proto.String(operation),
				RequestMethod: proto.String(r.Method),
				RequestUrl:    proto.String(r.URL.Path),
				Duration:      proto.Int64(duration),
				OperatorIp:    proto.String(httpx.GetRemoteAddr(r)),
				OperatorId:    proto.Int64(mixins.GetCurrentUserId(ctx)),
				OperatorName:  proto.String(mixins.GetCurrentUserName(ctx)),
				TenantId:      proto.Int64(mixins.GetCurrentTenantId(ctx)),
			}

			// 异步写入，失败仅记录日志，不影响主流程。使用独立 context 避免请求结束后被取消。
			go func() {
				bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if _, err := cli.CreateOperationLog(bgCtx, log); err != nil {
					logx.Errorf("write operation log error: %v", err)
				}
			}()
		}
	}
}

// parseModuleOperation 从请求路径推导模块名与操作名。
// 例：/system/user/create → ("user", "create")；/system/dict/data/list → ("dict", "list")
func parseModuleOperation(path string) (module, operation string) {
	segs := make([]string, 0, 4)
	for _, s := range strings.Split(path, "/") {
		if s != "" {
			segs = append(segs, s)
		}
	}
	switch len(segs) {
	case 0:
		return "", ""
	case 1:
		return segs[0], ""
	default:
		// 第 2 段作为模块（system 之后），末段作为操作
		return segs[1], segs[len(segs)-1]
	}
}
