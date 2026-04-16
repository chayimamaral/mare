package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"

	webview "github.com/webview/webview_go"
)

// Agora o padrão é local, sem os ".."
//
//go:embed all:out
var assets embed.FS

func main() {
	// Pega o conteúdo de dentro da pasta 'out'
	content, err := fs.Sub(assets, "out")
	if err != nil {
		log.Fatal("Erro ao acessar pasta out:", err)
	}

	go func() {
		log.Println("Servidor local rodando em http://localhost:9000")
		fsys := http.FS(content)
		fileServer := http.FileServer(fsys)

		// SPA fallback: se a rota não existir como arquivo/pasta exportada, devolve /index.html
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := r.URL.Path
			if p == "" || p == "/" {
				fileServer.ServeHTTP(w, r)
				return
			}

			clean := path.Clean("/" + p)
			clean = strings.TrimPrefix(clean, "/")

			// tenta servir o arquivo/pasta real exportado
			if f, err := fsys.Open(clean); err == nil {
				_ = f.Close()
				fileServer.ServeHTTP(w, r)
				return
			}

			// fallback para SPA
			r2 := r.Clone(r.Context())
			r2.URL.Path = "/index.html"
			fileServer.ServeHTTP(w, r2)
		})

		if err := http.ListenAndServe(":9000", handler); err != nil {
			log.Fatal(err)
		}
	}()

	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("VContab Desktop")
	w.SetSize(1280, 800, webview.HintNone)
	w.Navigate("http://localhost:9000")
	w.Run()
}
