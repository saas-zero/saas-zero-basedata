package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/zeromicro/go-zero/core/logx"

	casbinapi "github.com/casbin/casbin/v2"
)

func getRoleCodesFromCtx(ctx context.Context) []string {
	return GetRoleCodes(ctx)
}

// CasbinAuth returns HTTP middleware enforcing Casbin Domain RBAC.
// RoleCodes are read from JWT claims (set by JwtAuth middleware via context),
// then checked against Casbin policy for each role.
// fail-closed：disabled=false 时 enforcer 为 nil（初始化失败）一律 503 拒绝，
// 绝不放行；仅显式 enabled=false（本地开发配置 casbinDisabled）才放行。
func CasbinAuth(enf *casbinapi.SyncedEnforcer, disabled bool) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if disabled {
				next(w, r)
				return
			}
			if enf == nil {
				writeJSON(w, http.StatusInternalServerError, errno.AuthServiceUnavailable)
				return
			}
			if r.URL.Path == "" || r.URL.Path[0] != '/' {
				next(w, r)
				return
			}
			if len(r.URL.Path) >= 6 && r.URL.Path[:6] == "/init/" {
				next(w, r)
				return
			}
			// /system/api/mine 只返回"当前登录用户自己拥有的 API"（自带租户/user 隔离），
			// 且它是"分配 API"弹窗的数据源，不应受角色 Casbin 策略限制（新角色还没策略时也能打开）。
			if r.URL.Path == "/system/api/mine" {
				next(w, r)
				return
			}
			tenantId := mixins.GetCurrentTenantId(r.Context())

			// Read role codes from JWT claims (set by jwtauth middleware)
			roleCodes := getRoleCodesFromCtx(r.Context())
			if len(roleCodes) == 0 {
				writeJSON(w, http.StatusForbidden, errno.NoRoles)
				return
			}
			path := r.URL.Path
			method := r.Method
			dom := strconv.FormatInt(tenantId, 10)
			allowed := false
			for _, roleCode := range roleCodes {
				ok, err := enf.Enforce(roleCode, dom, path, method)
				if err != nil {
					// fail-closed：Casbin 引擎异常时拒绝访问，避免绕过权限校验
					logx.Errorf("Casbin enforce error: role=%s, dom=%s, path=%s, method=%s, err=%v",
						roleCode, dom, path, method, err)
					writeJSON(w, http.StatusInternalServerError, errno.AuthServiceUnavailable)
					return
				}
				if ok {
					allowed = true
					break
				}
			}
			if !allowed {
				writeJSON(w, http.StatusForbidden, errno.ForbiddenOperation)
				return
			}
			next(w, r)
		}
	}
}
