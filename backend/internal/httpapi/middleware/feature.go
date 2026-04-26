package middleware

import (
	"net/http"
	"strings"
)

// RequireFeature valida o slug no JWT (sem consulta ao banco). SUPER ignora a checagem.
// Tokens antigos sem feature_slugs: liberado (compatibilidade até novo login).
func RequireFeature(featureSlug string) func(http.Handler) http.Handler {
	want := strings.ToLower(strings.TrimSpace(featureSlug))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role := strings.ToUpper(strings.TrimSpace(Role(r.Context())))
			if role == "SUPER" {
				next.ServeHTTP(w, r)
				return
			}

			slugs := FeatureSlugs(r.Context())
			if slugs == nil {
				next.ServeHTTP(w, r)
				return
			}

			for _, s := range slugs {
				if strings.ToLower(strings.TrimSpace(s)) == want {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
