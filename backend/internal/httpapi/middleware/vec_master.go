package middleware

import (
	"net/http"
	"strings"
)

// RequireTenantVecMaster restringe a usuários do tenant master (VEC Sistemas), perfis ADMIN ou SUPER.
func RequireTenantVecMaster() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := strings.ToUpper(strings.TrimSpace(Role(r.Context())))
			if role != "ADMIN" && role != "SUPER" {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			if !TenantIsVecMaster(r.Context()) {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
