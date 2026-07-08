package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/errno"
	"github.com/saas-zero/saas-zero-common/pkg/jwt"
	"github.com/saas-zero/saas-zero-common/pkg/redis"
)

type ctxKey string

const roleCodesKey ctxKey = "role_codes"

func GetRoleCodes(ctx context.Context) []string {
	if v, ok := ctx.Value(roleCodesKey).([]string); ok {
		return v
	}
	return nil
}

func JwtAuth(secret string, rds *redis.Client) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/init/") {
				next(w, r)
				return
			}
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, errno.MissingAuthHeader.JSON(), http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, errno.InvalidAuthHeader.JSON(), http.StatusUnauthorized)
				return
			}
			claims, err := jwt.Parse(parts[1], secret)
			if err != nil {
				http.Error(w, errno.InvalidToken.JSON(), http.StatusUnauthorized)
				return
			}
			// Verify token exists in Redis (allows remote invalidation)
			if claims.ID != "" {
				exists, err := rds.Exists(fmt.Sprintf("token:%s", claims.ID))
				if err != nil || !exists {
					http.Error(w, errno.TokenInvalidated.JSON(), http.StatusUnauthorized)
					return
				}
			}
			// Verify token version matches Redis (invalidates on role/permission changes)
			if claims.TokenVersion > 0 {
				tv, err := rds.Get(fmt.Sprintf("token_version:%d", claims.UserId))
				if err != nil || tv == "" || tv != fmt.Sprintf("%d", claims.TokenVersion) {
					http.Error(w, errno.TokenVersionMismatch.JSON(), http.StatusUnauthorized)
					return
				}
			}
			ctx := r.Context()
			// Inject JWT claims into request context for downstream use.
			// - UserId/UserName/TenantId are used by ent mixin hooks for audit fields.
			// - RoleCodes are consumed by CasbinAuth middleware for authorization checks,
			//   eliminating a gRPC call to GetUserRoleCodes on every request.
			ctx = mixins.SetCurrentUserId(ctx, claims.UserId)
			ctx = mixins.SetCurrentUserName(ctx, claims.UserName)
			ctx = mixins.SetCurrentTenantId(ctx, claims.TenantId)
			ctx = context.WithValue(ctx, roleCodesKey, claims.RoleCodes)
			next(w, r.WithContext(ctx))
		}
	}
}
