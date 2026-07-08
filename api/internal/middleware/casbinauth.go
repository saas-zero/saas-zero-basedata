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
// If casbin is nil (graceful degradation), all requests pass through.
func CasbinAuth(enf *casbinapi.SyncedEnforcer) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if enf == nil {
				next(w, r)
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
			tenantId := mixins.GetCurrentTenantId(r.Context())

			// Read role codes from JWT claims (set by jwtauth middleware)
			roleCodes := getRoleCodesFromCtx(r.Context())
			if len(roleCodes) == 0 {
				http.Error(w, errno.NoRoles.JSON(), http.StatusForbidden)
				return
			}
			path := r.URL.Path
			method := r.Method
			dom := strconv.FormatInt(tenantId, 10)
			allowed := false
			for _, roleCode := range roleCodes {
				ok, err := enf.Enforce(roleCode, dom, path, method)
				if err != nil {
					logx.Errorf("Casbin enforce error: role=%s, dom=%s, path=%s, method=%s, err=%v",
						roleCode, dom, path, method, err)
					allowed = true
					break
				}
				if ok {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, errno.ForbiddenOperation.JSON(), http.StatusForbidden)
				return
			}
			next(w, r)
		}
	}
}
