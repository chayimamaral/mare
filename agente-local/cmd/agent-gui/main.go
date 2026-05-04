//go:build gui
// +build gui

// Interface com Fyne (Windows/Linux). Makefile: CGO + cross MinGW para .exe.

package main

import (
	"context"
	_ "embed"
	"fmt"
	"image"
	"log"
	"os/signal"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/chayimamaral/vecx/agente-local/internal/config"
	"github.com/chayimamaral/vecx/agente-local/internal/httpserver"
	"github.com/chayimamaral/vecx/agente-local/internal/images"
	"github.com/chayimamaral/vecx/agente-local/internal/provider/pkcs11"
	"github.com/chayimamaral/vecx/agente-local/internal/settings"
	"github.com/chayimamaral/vecx/agente-local/internal/usecase"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

//go:embed build_version.txt
var embeddedBuildVersion []byte

func resolvedVersion() string {
	v := strings.TrimSpace(string(embeddedBuildVersion))
	if v != "" {
		return v
	}
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}
	var rev string
	for _, s := range bi.Settings {
		if s.Key == "vcs.revision" {
			rev = s.Value
			break
		}
	}
	if rev != "" {
		if len(rev) > 12 {
			rev = rev[:12]
		}
		return "dev+" + rev
	}
	if bi.Main.Version != "" && bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}
	return "dev"
}

type statusState int

const (
	statusIdle statusState = iota
	statusDetected
	statusNotDetected
)

type guiState struct {
	mu sync.Mutex

	logLbl    *widget.Label
	statusDot *widget.RichText

	logs   []string
	status statusState
}

const logLineTimeFormat = "2006/01/02 15:04:05"

func (s *guiState) appendLog(msg string) {
	s.mu.Lock()
	line := time.Now().Format(logLineTimeFormat) + " " + msg
	s.logs = append(s.logs, line)
	if len(s.logs) > 500 {
		s.logs = s.logs[len(s.logs)-500:]
	}
	text := strings.Join(s.logs, "\n")
	logLbl := s.logLbl
	s.mu.Unlock()

	if logLbl != nil && fyne.CurrentApp() != nil {
		fyne.Do(func() {
			logLbl.SetText(text)
		})
	}
}

func statusColorName(st statusState) fyne.ThemeColorName {
	switch st {
	case statusDetected:
		return theme.ColorNameSuccess
	case statusNotDetected:
		return theme.ColorNameError
	default:
		return theme.ColorNameDisabled
	}
}

func (s *guiState) setStatus(st statusState) {
	s.mu.Lock()
	s.status = st
	rt := s.statusDot
	s.mu.Unlock()

	if rt != nil && fyne.CurrentApp() != nil {
		fyne.Do(func() {
			rt.Segments = []widget.RichTextSegment{
				&widget.TextSegment{
					Text: " ● ",
					Style: widget.RichTextStyle{
						ColorName: statusColorName(st),
						Inline:    true,
						SizeName:  theme.SizeNameHeadingText,
					},
				},
			}
			rt.Refresh()
		})
	}
}

func rasterizeSVG(svg string, w, h int) image.Image {
	icon, err := oksvg.ReadIconStream(strings.NewReader(svg))
	if err != nil {
		return image.NewRGBA(image.Rect(0, 0, w, h))
	}
	icon.SetTarget(0, 0, float64(w), float64(h))
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	scanner := rasterx.NewScannerGV(w, h, img, img.Bounds())
	raster := rasterx.NewDasher(w, h, scanner)
	icon.Draw(raster, 1.0)
	return img
}

func appWindowTitle() string {
	return "Agente VECX - " + resolvedVersion()
}

func main() {
	cfg := config.Load()
	provider := pkcs11.NewProvider(cfg.PKCS11LibraryLinux, cfg.PKCS11LibraryWindow)
	store, errStore := settings.DefaultStore()
	if errStore != nil {
		log.Printf("configuracao local: %v", errStore)
	}
	signUC := usecase.NewSignUseCase(provider, store)
	state := &guiState{status: statusIdle}

	handler := httpserver.NewHandler(signUC, state.appendLog)
	server := httpserver.NewServer(cfg.HTTPAddr, cfg.AllowedOrigins, cfg.SharedSecret, handler)

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.NewWithID("com.vecsistemas.vecx.agent")

	w := a.NewWindow(appWindowTitle())
	w.Resize(fyne.NewSize(820, 640))
	w.SetFixedSize(false)

	w.SetCloseIntercept(func() {
		stop()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_ = server.Shutdown(shutdownCtx)
		cancel()
		w.Close()
	})

	go func() {
		<-stopCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		a.Quit()
	}()

	logo := canvas.NewImageFromImage(rasterizeSVG(images.VecxLogoSVG, 180, 64))
	logo.SetMinSize(fyne.NewSize(180, 64))
	logo.FillMode = canvas.ImageFillContain

	headTitle := widget.NewLabel("Agente VECX")
	headTitle.TextStyle = fyne.TextStyle{Bold: true}
	head := container.NewHBox(logo, headTitle)

	verLbl := widget.NewLabel("Versao " + resolvedVersion())
	hintA1 := widget.NewLabel("Certificados A1 em disco: subpastas cert_clientes/ e cert_contador/; arquivos {CNPJ|CPF}.pfx")
	hintA1.Wrapping = fyne.TextWrapWord

	rootEntry := widget.NewEntry()
	rootEntry.SetPlaceHolder("Pasta raiz (ex.: C:\\Users\\...\\certs-vecx)")

	preferA3 := widget.NewCheck("Preferir A3 (ignorar .pfx locais e usar token)", nil)
	if store != nil {
		st := store.Load()
		rootEntry.SetText(st.CertRootDir)
		preferA3.SetChecked(st.PreferA3)
	}

	browseBtn := widget.NewButton("Procurar...", func() {
		dialog.ShowFolderOpen(func(u fyne.ListableURI, err error) {
			if err != nil {
				state.appendLog("Seletor de pasta: " + err.Error())
				return
			}
			if u == nil {
				return
			}
			rootEntry.SetText(u.Path())
		}, w)
	})

	saveBtn := widget.NewButton("Salvar", func() {
		if store == nil {
			state.appendLog("Configuracao local indisponivel")
			return
		}
		st := settings.AgentSettings{
			CertRootDir: strings.TrimSpace(rootEntry.Text),
			PreferA3:    preferA3.Checked,
		}
		if err := store.Save(st); err != nil {
			state.appendLog("Falha ao salvar configuracao: " + err.Error())
		} else {
			state.appendLog("Configuracao salva (pasta raiz / preferir A3)")
		}
	})

	pathRow := container.NewBorder(nil, nil, nil, container.NewHBox(browseBtn, saveBtn), rootEntry)

	statusMark := widget.NewRichText(&widget.TextSegment{
		Text: " ● ",
		Style: widget.RichTextStyle{
			ColorName: statusColorName(statusIdle),
			Inline:    true,
			SizeName:  theme.SizeNameHeadingText,
		},
	})
	state.statusDot = statusMark

	detectBtn := widget.NewButton("Detectar", func() {
		go func() {
			state.appendLog("Acao manual: Detectar certificado A3")
			certs, err := signUC.ListCertificates(context.Background())
			if err != nil {
				state.setStatus(statusNotDetected)
				state.appendLog("Deteccao falhou: " + err.Error())
				return
			}
			if len(certs) == 0 {
				state.setStatus(statusNotDetected)
				state.appendLog("Nenhum certificado A3 detectado")
			} else {
				state.setStatus(statusDetected)
				state.appendLog(fmt.Sprintf("%d certificado(s) detectado(s)", len(certs)))
			}
		}()
	})

	legend := widget.NewLabel("Cinza=aguardando | Verde=detectado | Vermelho=nao detectado")
	detectRow := container.NewHBox(detectBtn, statusMark, legend)

	logHint := widget.NewLabel("Log: selecione texto e Ctrl+C para copiar; roda do mouse para rolar")
	logLbl := widget.NewLabel("")
	logLbl.Wrapping = fyne.TextWrapWord
	logLbl.TextStyle = fyne.TextStyle{Monospace: true}
	logLbl.Selectable = true
	state.logLbl = logLbl
	scrollLog := container.NewScroll(logLbl)
	scrollLog.SetMinSize(fyne.NewSize(760, 260))

	inner := container.NewVBox(
		head,
		widget.NewSeparator(),
		verLbl,
		hintA1,
		pathRow,
		preferA3,
		detectRow,
		logHint,
		scrollLog,
	)
	w.SetContent(container.NewPadded(inner))

	state.appendLog("versao " + resolvedVersion())
	state.appendLog(fmt.Sprintf("agente local iniciado em http://%s", cfg.HTTPAddr))
	if len(cfg.AllowedOrigins) > 0 {
		state.appendLog(fmt.Sprintf("cors liberado para: %v", cfg.AllowedOrigins))
	}
	if cfg.SharedSecret != "" {
		state.appendLog("autenticacao local habilitada via X-Local-Agent-Secret")
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !strings.Contains(strings.ToLower(err.Error()), "closed") {
			state.appendLog("Erro no servidor: " + err.Error())
		}
	}()

	w.ShowAndRun()
}
