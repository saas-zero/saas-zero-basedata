package middleware

import (
	"net/http"
	"strconv"

	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"

	casbinapi "github.com/casbin/casbin/v2"
)

func CasbinAuth(enf *casbinapi.SyncedEnforcer) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "" || r.URL.Path[0] != '/' {
				next(w, r)
				return
			}
			if r.URL.Path[:6] == "/init/" {
				next(w, r)
				return
			}
			tenantId := mixins.GetCurrentTenantId(r.Context())
			roleCodes := GetRoleCodes(r.Context())
			if len(roleCodes) == 0 {
				http.Error(w, `{"code":403,"msg":"no roles"}`, http.StatusForbidden)
				return
			}
			path := r.URL.Path
			method := r.Method
			dom := strconv.FormatInt(tenantId, 10)
			allowed := false
			for _, roleCode := range roleCodes {
				if ok, _ := enf.Enforce(roleCode, dom, path, method); ok {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, `{"code":403,"msg":"forbidden"}`, http.StatusForbidden)
				return
			}
			next(w, r)
		}
	}
}
