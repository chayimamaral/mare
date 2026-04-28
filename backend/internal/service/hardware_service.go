package service

import (
	"context"
	"os/exec"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

type HardwareDevice struct {
	DeviceID string `json:"device_id"`
	Path     string `json:"path"`
	Status   string `json:"status"`
}

type HardwareService struct{}

func NewHardwareService() *HardwareService {
	return &HardwareService{}
}

// ListLocalDevices lista dispositivos locais de forma best-effort no host do backend.
// Em runtime desktop/local, reflete o hardware da máquina executando o binário.
func (s *HardwareService) ListLocalDevices(ctx context.Context) ([]HardwareDevice, error) {
	if runtime.GOOS == "windows" {
		return s.listLocalDevicesWindows(ctx)
	}
	return s.listLocalDevicesLinux(ctx)
}

func (s *HardwareService) listLocalDevicesLinux(ctx context.Context) ([]HardwareDevice, error) {
	_ = ctx
	devs := map[string]HardwareDevice{}

	// 1) Barramento USB geral (pendrive, token, leitor, etc).
	// Sempre registra o nó USB bruto para nao perder dispositivo sem metadata.
	if entries, err := os.ReadDir("/sys/bus/usb/devices"); err == nil {
		for _, e := range entries {
			name := strings.TrimSpace(e.Name())
			if name == "" || strings.HasPrefix(name, "usb") || strings.Contains(name, ":") {
				continue
			}

			base := filepath.Join("/sys/bus/usb/devices", name)
			man := readSysText(base, "manufacturer")
			prod := readSysText(base, "product")
			serial := readSysText(base, "serial")
			vid := readSysText(base, "idVendor")
			pid := readSysText(base, "idProduct")
			bus := readSysText(base, "busnum")
			devnum := readSysText(base, "devnum")
			uevent := readSysText(base, "uevent")
			// Mesmo sem metadados legiveis, mantem dispositivo usb_device.
			isUSBDevice := strings.Contains(uevent, "DEVTYPE=usb_device")

			label := strings.TrimSpace(strings.TrimSpace(man + " " + prod))
			if label == "" {
				label = "USB device"
			}
			if serial != "" {
				label += " | serial: " + serial
			}
			if vid != "" || pid != "" {
				label += " | vid:pid " + strings.TrimSpace(vid+":"+pid)
			}
			if bus != "" || devnum != "" {
				label += " | bus/dev " + strings.TrimSpace(bus+"/"+devnum)
			}
			if isUSBDevice {
				label += " | usb_device"
			}
			if vid == "" && pid == "" && man == "" && prod == "" && !isUSBDevice {
				label += " | sem_metadata"
			}

			id := "usb-" + name
			devs[id] = HardwareDevice{
				DeviceID: id,
				Path:     "/sys/bus/usb/devices/" + name + " | " + label,
				Status:   "Conectado",
			}
		}
	}

	// 2) Blocos em /sys/block (Linux/Fedora).
	if entries, err := os.ReadDir("/sys/block"); err == nil {
		for _, e := range entries {
			name := strings.TrimSpace(e.Name())
			if name == "" {
				continue
			}
			id := "disk-" + name
			path := "/dev/" + name
			status := "Conectado"
			if strings.HasPrefix(name, "loop") || strings.HasPrefix(name, "ram") {
				status = "Virtual"
			}
			devs[id] = HardwareDevice{
				DeviceID: id,
				Path:     path,
				Status:   status,
			}
		}
	}

	// 3) Links de USB em /dev/disk/by-id (quando disponível).
	if entries, err := os.ReadDir("/dev/disk/by-id"); err == nil {
		for _, e := range entries {
			n := strings.TrimSpace(e.Name())
			if n == "" {
				continue
			}
			if !strings.Contains(strings.ToLower(n), "usb") {
				continue
			}
			linkPath := filepath.Join("/dev/disk/by-id", n)
			target, err := filepath.EvalSymlinks(linkPath)
			if err != nil {
				target = linkPath
			}
			id := "usb-" + n
			devs[id] = HardwareDevice{
				DeviceID: id,
				Path:     target,
				Status:   "Conectado",
			}
		}
	}

	// 4) Fallback em /dev para tokens/serial USB (ttyACM, ttyUSB, hidraw).
	if entries, err := os.ReadDir("/dev"); err == nil {
		for _, e := range entries {
			n := strings.TrimSpace(e.Name())
			if n == "" {
				continue
			}
			if !(strings.HasPrefix(n, "ttyACM") || strings.HasPrefix(n, "ttyUSB") || strings.HasPrefix(n, "hidraw")) {
				continue
			}
			id := "dev-" + n
			devs[id] = HardwareDevice{
				DeviceID: id,
				Path:     "/dev/" + n,
				Status:   "Conectado",
			}
		}
	}

	// 5) lsusb complementar (sempre): adiciona descricao do barramento.
	lines := runCommandLines(ctx, 4*time.Second, "lsusb")
	for _, ln := range lines {
		t := strings.TrimSpace(ln)
		if t == "" {
			continue
		}
		// Ex: Bus 003 Device 005: ID 058f:6387 Generic Flash Disk
		id := "lsusb-" + sanitizeID(t)
		devs[id] = HardwareDevice{
			DeviceID: id,
			Path:     t,
			Status:   "Conectado",
		}
	}

	out := make([]HardwareDevice, 0, len(devs))
	for _, d := range devs {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DeviceID < out[j].DeviceID })
	return out, nil
}

func (s *HardwareService) listLocalDevicesWindows(ctx context.Context) ([]HardwareDevice, error) {
	devs := map[string]HardwareDevice{}

	// 1) USB Plug and Play devices (token, pendrive, leitora, etc)
	psUSB := `$ErrorActionPreference='SilentlyContinue'; Get-PnpDevice -PresentOnly | Where-Object { $_.InstanceId -like 'USB*' } | ForEach-Object { "$($_.InstanceId)|$($_.FriendlyName)" }`
	if lines, err := runPowerShellLines(ctx, psUSB); err == nil {
		for _, ln := range lines {
			parts := strings.SplitN(ln, "|", 2)
			id := strings.TrimSpace(parts[0])
			if id == "" {
				continue
			}
			name := id
			if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
				name = strings.TrimSpace(parts[1])
			}
			devs["usb-"+id] = HardwareDevice{
				DeviceID: "usb-" + id,
				Path:     name,
				Status:   "Conectado",
			}
		}
	}

	// 2) USB disk drives
	psDisk := `$ErrorActionPreference='SilentlyContinue'; Get-CimInstance Win32_DiskDrive | Where-Object { $_.InterfaceType -eq 'USB' } | ForEach-Object { "$($_.DeviceID)|$($_.Model)" }`
	if lines, err := runPowerShellLines(ctx, psDisk); err == nil {
		for _, ln := range lines {
			parts := strings.SplitN(ln, "|", 2)
			id := strings.TrimSpace(parts[0])
			if id == "" {
				continue
			}
			name := id
			if len(parts) > 1 && strings.TrimSpace(parts[1]) != "" {
				name = strings.TrimSpace(parts[1])
			}
			devs["disk-"+id] = HardwareDevice{
				DeviceID: "disk-" + id,
				Path:     name,
				Status:   "Conectado",
			}
		}
	}

	out := make([]HardwareDevice, 0, len(devs))
	for _, d := range devs {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].DeviceID < out[j].DeviceID })
	return out, nil
}

func readSysText(base, file string) string {
	p := filepath.Join(base, file)
	b, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func runPowerShellLines(ctx context.Context, script string) ([]string, error) {
	cctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	rawLines := strings.Split(string(out), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, ln := range rawLines {
		t := strings.TrimSpace(ln)
		if t == "" {
			continue
		}
		lines = append(lines, t)
	}
	return lines, nil
}

func runCommandLines(ctx context.Context, timeout time.Duration, name string, args ...string) []string {
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, name, args...)
	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	rawLines := strings.Split(string(out), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, ln := range rawLines {
		t := strings.TrimSpace(ln)
		if t != "" {
			lines = append(lines, t)
		}
	}
	return lines
}

func sanitizeID(v string) string {
	s := strings.ToLower(strings.TrimSpace(v))
	repl := []string{" ", "-", ":", "/", ".", ","}
	for _, r := range repl {
		s = strings.ReplaceAll(s, r, "_")
	}
	return s
}

