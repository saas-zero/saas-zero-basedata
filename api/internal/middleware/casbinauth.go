package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/saas-zero/saas-zero-basedata/rpc/apps"
	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"

	casbinapi "github.com/casbin/casbin/v2"
)

type roleCodesKeyType string

const roleCodesCtxKey roleCodesKeyType = "role_codes"

func withRoleCodes(ctx context.Context, codes []string) context.Context {
	return context.WithValue(ctx, roleCodesCtxKey, codes)
}

func getRoleCodesFromCtx(ctx context.Context) []string {
	if v, ok := ctx.Value(roleCodesCtxKey).([]string); ok {
		return v
	}
	return nil
}

func CasbinAuth(enf *casbinapi.SyncedEnforcer, usersClient apps.SysUsersClient) func(http.HandlerFunc) http.HandlerFunc {
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
			userId := mixins.GetCurrentUserId(r.Context())

			// Fetch role codes via gRPC (not from JWT)
			roleCodes := getRoleCodesFromCtx(r.Context())
			if roleCodes == nil {
				resp, err := usersClient.GetUserRoleCodes(r.Context(), &apps.IdReq{Id: userId})
				if err != nil || resp == nil {
					http.Error(w, errno.InvalidToken.JSON(), http.StatusUnauthorized)
					return
				}
				roleCodes = resp.GetCodes()
				r = r.WithContext(withRoleCodes(r.Context(), roleCodes))
			}

			if len(roleCodes) == 0 {
				http.Error(w, errno.NoRoles.JSON(), http.StatusForbidden)
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
				http.Error(w, errno.ForbiddenOperation.JSON(), http.StatusForbidden)
				return
			}
			next(w, r)
		}
	}
}
