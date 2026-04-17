package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecontab/backend/internal/auth"
	"github.com/chayimamaral/vecontab/backend/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

type TenantSchemaResolver func(ctx context.Context, tenantID string) (string, error)

const (
	userIDKey   contextKey = "userID"
	roleKey     contextKey = "role"
	tenantIDKey contextKey = "tenantID"
	tenantSchemaKey contextKey = "tenantSchema"
)

var tenantSchemaResolver TenantSchemaResolver
var tenantConnPool *pgxpool.Pool

func SetTenantSchemaResolver(resolver TenantSchemaResolver) {
	tenantSchemaResolver = resolver
}

func SetTenantConnPool(pool *pgxpool.Pool) {
	tenantConnPool = pool
}

func RequireAuth(tokens *auth.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, err := tokens.Parse(parts[1])
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			ctx = context.WithValue(ctx, tenantIDKey, claims.Tenant.ID)
			tenantSchema := strings.TrimSpace(claims.Tenant.SchemaName)
			if tenantSchema == "" && tenantSchemaResolver != nil && strings.TrimSpace(claims.Tenant.ID) != "" {
				if resolved, err := tenantSchemaResolver(r.Context(), claims.Tenant.ID); err == nil {
					tenantSchema = strings.TrimSpace(resolved)
				}
			}
			ctx = context.WithValue(ctx, tenantSchemaKey, tenantSchema)
			if tenantConnPool != nil && tenantSchema != "" {
				conn, err := db.AcquireTenantConn(ctx, tenantConnPool, tenantSchema)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
					return
				}
				defer conn.Release()
				ctx = db.ContextWithConn(ctx, conn)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireAnyRole(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		role = strings.TrimSpace(strings.ToUpper(role))
		if role != "" {
			allowed[role] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := strings.TrimSpace(strings.ToUpper(Role(r.Context())))
			if role == "" {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			if _, ok := allowed[role]; !ok {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func UserID(ctx context.Context) string {
	value, _ := ctx.Value(userIDKey).(string)
	return value
}

func Role(ctx context.Context) string {
	value, _ := ctx.Value(roleKey).(string)
	return value
}

func TenantID(ctx context.Context) string {
	value, _ := ctx.Value(tenantIDKey).(string)
	return value
}

func TenantSchema(ctx context.Context) string {
	value, _ := ctx.Value(tenantSchemaKey).(string)
	return value
}
