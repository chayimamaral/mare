//go:build gui
// +build gui

package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/chayimamaral/vecx/agente-local/internal/config"
	"github.com/chayimamaral/vecx/agente-local/internal/httpserver"
	"github.com/chayimamaral/vecx/agente-local/internal/images"
	"github.com/chayimamaral/vecx/agente-local/internal/provider/pkcs11"
	"github.com/chayimamaral/vecx/agente-local/internal/usecase"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

type statusState int

const (
	statusIdle statusState = iota
	statusDetected
	statusNotDetected
)

type guiState struct {
	mu         sync.Mutex
	invalidate func()
	logs       []string
	status     statusState
	detectBt   widget.Clickable
}

// Mesmo estilo visual do pacote log do Go (LstdFlags: data, hora, mensagem).
const logLineTimeFormat = "2006/01/02 15:04:05"

func (s *guiState) setInvalidate(f func()) {
	s.mu.Lock()
	s.invalidate = f
	s.mu.Unlock()
}

func (s *guiState) appendLog(msg string) {
	s.mu.Lock()
	line := time.Now().Format(logLineTimeFormat) + " " + msg
	s.logs = append(s.logs, line)
	if len(s.logs) > 500 {
		s.logs = s.logs[len(s.logs)-500:]
	}
	inv := s.invalidate
	s.mu.Unlock()
	if inv != nil {
		inv()
	}
}

func (s *guiState) setStatus(st statusState) {
	s.mu.Lock()
	s.status = st
	s.mu.Unlock()
}

func (s *guiState) snapshot() (statusState, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status, strings.Join(s.logs, "\n")
}

func (s *guiState) snapshotLines() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, len(s.logs))
	copy(out, s.logs)
	return out
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

func drawStatusDot(gtx layout.Context, st statusState) layout.Dimensions {
	col := color.NRGBA{R: 140, G: 140, B: 140, A: 255}
	switch st {
	case statusDetected:
		col = color.NRGBA{R: 33, G: 160, B: 70, A: 255}
	case statusNotDetected:
		col = color.NRGBA{R: 190, G: 48, B: 48, A: 255}
	}
	size := gtx.Dp(unit.Dp(14))
	defer clip.Ellipse{Max: image.Pt(size, size)}.Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, col)
	return layout.Dimensions{Size: image.Pt(size, size)}
}

func main() {
	cfg := config.Load()
	provider := pkcs11.NewProvider(cfg.PKCS11LibraryLinux, cfg.PKCS11LibraryWindow)
	signUC := usecase.NewSignUseCase(provider)
	state := &guiState{status: statusIdle}

	handler := httpserver.NewHandler(signUC, state.appendLog)
	server := httpserver.NewServer(cfg.HTTPAddr, cfg.AllowedOrigins, cfg.SharedSecret, handler)

	stopCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-stopCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
		os.Exit(0)
	}()

	go func() {
		w := new(app.Window)
		state.setInvalidate(w.Invalidate)
		w.Option(
			app.Title("Agente VECX"),
			app.Size(unit.Dp(820), unit.Dp(640)),
		)

		// Depois de Invalidate estar ligado, para o textbox atualizar de imediato.
		// Textos iguais ao vecx-agent-cli (cmd/agent/main.go), com prefixo data/hora estilo log.
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

		th := material.NewTheme()
		th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

		var ops op.Ops
		var logList widget.List
		logList.Axis = layout.Vertical
		logList.ScrollToEnd = true
		logo := paint.NewImageOp(rasterizeSVG(images.VecxLogoSVG, 180, 64))

		for {
			e := w.Event()
			switch e := e.(type) {
			case app.DestroyEvent:
				if e.Err != nil {
					log.Printf("janela encerrada com erro: %v", e.Err)
				}
				stop()
				return
			case app.FrameEvent:
				ops.Reset()
				gtx := app.NewContext(&ops, e)

				for state.detectBt.Clicked(gtx) {
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
				}

				st, _ := state.snapshot()
				lines := state.snapshotLines()

				layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Max.Y = min(gtx.Constraints.Max.Y, gtx.Dp(unit.Dp(80)))
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									img := widget.Image{Src: logo, Fit: widget.Contain}
									return img.Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
								layout.Rigid(material.H4(th, "Agente VECX").Layout),
							)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(12)}.Layout),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return material.Button(th, &state.detectBt, "Detectar").Layout(gtx)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(12)}.Layout),
								layout.Rigid(func(gtx layout.Context) layout.Dimensions {
									return drawStatusDot(gtx, st)
								}),
								layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout),
								layout.Rigid(material.Body2(th, "Cinza=aguardando | Verde=detectado | Vermelho=nao detectado").Layout),
							)
						}),
						layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
						layout.Rigid(material.Body2(th, "Mensagens (lista com rolagem; novas linhas seguem para o fim)").Layout),
						layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
						layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints = gtx.Constraints.AddMin(image.Pt(0, gtx.Dp(unit.Dp(260))))
							border := widget.Border{Color: color.NRGBA{R: 180, G: 180, B: 190, A: 255}, CornerRadius: unit.Dp(4), Width: unit.Dp(1)}
							return border.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
									return material.List(th, &logList).Layout(gtx, len(lines), func(gtx layout.Context, i int) layout.Dimensions {
										return layout.UniformInset(unit.Dp(2)).Layout(gtx, material.Body2(th, lines[i]).Layout)
									})
								})
							})
						}),
					)
				})

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}

