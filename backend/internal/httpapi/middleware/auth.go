package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chayimamaral/vecx/backend/internal/auth"
	"github.com/chayimamaral/vecx/backend/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

type TenantSchemaResolver func(ctx context.Context, tenantID string) (string, error)

const (
	userIDKey            contextKey = "userID"
	userNameKey          contextKey = "userName"
	roleKey              contextKey = "role"
	tenantIDKey          contextKey = "tenantID"
	tenantSchemaKey      contextKey = "tenantSchema"
	representativeIDKey  contextKey = "representativeID"
	featureSlugsKey      contextKey = "featureSlugs"
	tenantIsVecMasterKey contextKey = "tenantIsVecMaster"
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
			tokenStr := ""
			header := r.Header.Get("Authorization")
			if header != "" {
				parts := strings.SplitN(header, " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					tokenStr = parts[1]
				}
			}

			// Fallback para WebSocket / SSE
			if tokenStr == "" {
				tokenStr = r.URL.Query().Get("token")
			}

			if tokenStr == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, err := tokens.Parse(tokenStr)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.Subject)
			ctx = context.WithValue(ctx, userNameKey, claims.Nome)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			ctx = context.WithValue(ctx, tenantIDKey, claims.Tenant.ID)
			ctx = context.WithValue(ctx, representativeIDKey, strings.TrimSpace(claims.RepresentativeID))
			ctx = context.WithValue(ctx, featureSlugsKey, claims.FeatureSlugs)
			ctx = context.WithValue(ctx, tenantIsVecMasterKey, claims.Tenant.IsVecMaster)
			tenantSchema := strings.TrimSpace(claims.Tenant.SchemaName)
			if tenantSchema == "" && tenantSchemaResolver != nil && strings.TrimSpace(claims.Tenant.ID) != "" {
				if resolved, err := tenantSchemaResolver(r.Context(), claims.Tenant.ID); err == nil {
					tenantSchema = strings.TrimSpace(resolved)
				}
			}
			if strings.TrimSpace(claims.Tenant.ID) != "" && tenantSchema == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx = context.WithValue(ctx, tenantSchemaKey, tenantSchema)
			if tenantSchema != "" {
				if tenantConnPool == nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
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

func UserName(ctx context.Context) string {
	value, _ := ctx.Value(userNameKey).(string)
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

func RepresentativeID(ctx context.Context) string {
	value, _ := ctx.Value(representativeIDKey).(string)
	return value
}

func FeatureSlugs(ctx context.Context) []string {
	value, _ := ctx.Value(featureSlugsKey).([]string)
	return value
}

func TenantIsVecMaster(ctx context.Context) bool {
	value, _ := ctx.Value(tenantIsVecMasterKey).(bool)
	return value
}
