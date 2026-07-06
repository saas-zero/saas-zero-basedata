package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/saas-zero/saas-zero-common/pkg/ent/mixins"
	"github.com/saas-zero/saas-zero-common/pkg/jwt"
)

type ctxKey string

const roleCodesKey ctxKey = "role_codes"

func GetRoleCodes(ctx context.Context) []string {
	if v, ok := ctx.Value(roleCodesKey).([]string); ok {
		return v
	}
	return nil
}

func JwtAuth(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/init/") {
				next(w, r)
				return
			}
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, `{"code":401,"msg":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				http.Error(w, `{"code":401,"msg":"invalid authorization header"}`, http.StatusUnauthorized)
				return
			}
			claims, err := jwt.Parse(parts[1], secret)
			if err != nil {
				http.Error(w, `{"code":401,"msg":"invalid token"}`, http.StatusUnauthorized)
				return
			}
			ctx := r.Context()
			ctx = mixins.SetCurrentUserId(ctx, claims.UserId)
			ctx = mixins.SetCurrentUserName(ctx, claims.UserName)
			ctx = mixins.SetCurrentTenantId(ctx, claims.TenantId)
			ctx = context.WithValue(ctx, roleCodesKey, claims.RoleCodes)
			next(w, r.WithContext(ctx))
		}
	}
}
