package httpapi

import (
	"bytes"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

// ServeSPA entrega o export estático do Next; rotas sem arquivo físico recebem index.html.
// Não usa http.FileServer para a raiz nem para diretórios: o FileServer pode responder
// 301 com Location: ./ e gerar loop infinito atrás de proxies (Cloud Run).
func ServeSPA(root fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		p := normalizeWebPath(r.URL.Path)
		rel := strings.TrimPrefix(p, "/")

		if p == "/" || p == "/index.html" {
			serveEmbeddedFile(w, r, root, "index.html")
			return
		}

		fi, err := fs.Stat(root, rel)
		if err == nil && !fi.IsDir() {
			http.ServeFileFS(w, r, root, rel)
			return
		}

		serveEmbeddedFile(w, r, root, "index.html")
	})
}

func normalizeWebPath(p string) string {
	if p == "" {
		return "/"
	}
	p = path.Clean(p)
	if p == "." {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		return "/" + p
	}
	return p
}

func serveEmbeddedFile(w http.ResponseWriter, r *http.Request, fsys fs.FS, name string) {
	b, err := fs.ReadFile(fsys, name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(b))
}
