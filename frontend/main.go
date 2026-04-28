package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

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

	// Runtime JS leve para compatibilidade com listener estilo Wails:
	// window.runtime.EventsOn("hardware:usb-inserted", cb)
	w.Init(`
		(() => {
			if (!window.runtime) window.runtime = {};
			if (!window.__vecxEventMap) window.__vecxEventMap = {};
			if (!window.runtime.EventsOn) {
				window.runtime.EventsOn = function(eventName, cb) {
					if (!window.__vecxEventMap[eventName]) window.__vecxEventMap[eventName] = [];
					window.__vecxEventMap[eventName].push(cb);
					return function() {
						window.__vecxEventMap[eventName] = (window.__vecxEventMap[eventName] || []).filter(fn => fn !== cb);
					};
				};
			}
			if (!window.__vecxEmitEvent) {
				window.__vecxEmitEvent = function(eventName, payload) {
					const arr = window.__vecxEventMap[eventName] || [];
					arr.forEach(fn => {
						try { fn(payload); } catch (_) {}
					});
				};
			}
		})();
	`)
	w.Navigate("http://localhost:9000")

	go watchUSBAndEmit(w)
	w.Run()
}

func watchUSBAndEmit(w webview.WebView) {
	prev := map[string]struct{}{}
	seed := listUSBCandidates()
	for _, d := range seed {
		prev[d] = struct{}{}
	}

	tk := time.NewTicker(3 * time.Second)
	defer tk.Stop()
	for range tk.C {
		curList := listUSBCandidates()
		cur := map[string]struct{}{}
		for _, d := range curList {
			cur[d] = struct{}{}
		}

		inserted := make([]string, 0)
		for d := range cur {
			if _, ok := prev[d]; !ok {
				inserted = append(inserted, d)
			}
		}
		prev = cur

		if len(inserted) == 0 {
			continue
		}
		sort.Strings(inserted)
		payload := strings.Join(inserted, ",")
		w.Dispatch(func() {
			w.Eval(`window.__vecxEmitEvent && window.__vecxEmitEvent("hardware:usb-inserted", "` + jsEscape(payload) + `");`)
		})
	}
}

func listUSBCandidates() []string {
	if runtime.GOOS == "windows" {
		return listUSBCandidatesWindows()
	}
	return listUSBCandidatesLinux()
}

func listUSBCandidatesLinux() []string {
	out := make(map[string]struct{})

	// 1) Barramento USB geral: detecta token/leitor/perifericos nao-bloco.
	if entries, err := os.ReadDir("/sys/bus/usb/devices"); err == nil {
		for _, e := range entries {
			n := strings.TrimSpace(e.Name())
			if n == "" || strings.HasPrefix(n, "usb") || strings.Contains(n, ":") {
				continue
			}
			base := filepath.Join("/sys/bus/usb/devices", n)
			vid := readSysText(base, "idVendor")
			pid := readSysText(base, "idProduct")
			prod := readSysText(base, "product")
			man := readSysText(base, "manufacturer")
			if vid == "" && pid == "" && prod == "" && man == "" {
				continue
			}
			out["usb:"+n+":"+strings.TrimSpace(vid+":"+pid)] = struct{}{}
		}
	}

	// 2) /dev/disk/by-id com symlinks usb-*
	if entries, err := os.ReadDir("/dev/disk/by-id"); err == nil {
		for _, e := range entries {
			n := strings.TrimSpace(e.Name())
			if n == "" || !strings.Contains(strings.ToLower(n), "usb") {
				continue
			}
			p := filepath.Join("/dev/disk/by-id", n)
			if t, err := filepath.EvalSymlinks(p); err == nil {
				out[t] = struct{}{}
			} else {
				out[p] = struct{}{}
			}
		}
	}

	// 3) /sys/block/<dev>/removable == 1 (pendrive e similares)
	if entries, err := os.ReadDir("/sys/block"); err == nil {
		for _, e := range entries {
			dev := strings.TrimSpace(e.Name())
			if dev == "" {
				continue
			}
			remPath := filepath.Join("/sys/block", dev, "removable")
			b, err := os.ReadFile(remPath)
			if err != nil {
				continue
			}
			if strings.TrimSpace(string(b)) == "1" {
				out["/dev/"+dev] = struct{}{}
			}
		}
	}

	// 4) /dev nodes usuais para token/leitor
	if entries, err := os.ReadDir("/dev"); err == nil {
		for _, e := range entries {
			n := strings.TrimSpace(e.Name())
			if n == "" {
				continue
			}
			if strings.HasPrefix(n, "ttyACM") || strings.HasPrefix(n, "ttyUSB") || strings.HasPrefix(n, "hidraw") {
				out["/dev/"+n] = struct{}{}
			}
		}
	}

	arr := make([]string, 0, len(out))
	for d := range out {
		arr = append(arr, d)
	}
	sort.Strings(arr)
	return arr
}

func listUSBCandidatesWindows() []string {
	out := make(map[string]struct{})

	// 1) USB devices gerais
	psUSB := `$ErrorActionPreference='SilentlyContinue'; Get-PnpDevice -PresentOnly | Where-Object { $_.InstanceId -like 'USB*' } | ForEach-Object { "$($_.InstanceId)|$($_.FriendlyName)" }`
	for _, ln := range runPowerShellLines(psUSB) {
		if ln == "" {
			continue
		}
		out["usb:"+ln] = struct{}{}
	}

	// 2) USB disks
	psDisk := `$ErrorActionPreference='SilentlyContinue'; Get-CimInstance Win32_DiskDrive | Where-Object { $_.InterfaceType -eq 'USB' } | ForEach-Object { "$($_.DeviceID)|$($_.Model)" }`
	for _, ln := range runPowerShellLines(psDisk) {
		if ln == "" {
			continue
		}
		out["disk:"+ln] = struct{}{}
	}

	arr := make([]string, 0, len(out))
	for d := range out {
		arr = append(arr, d)
	}
	sort.Strings(arr)
	return arr
}

func runPowerShellLines(script string) []string {
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	raw := strings.Split(string(out), "\n")
	lines := make([]string, 0, len(raw))
	for _, ln := range raw {
		t := strings.TrimSpace(ln)
		if t != "" {
			lines = append(lines, t)
		}
	}
	return lines
}

func readSysText(base, file string) string {
	b, err := os.ReadFile(filepath.Join(base, file))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func jsEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", "")
	return s
}
