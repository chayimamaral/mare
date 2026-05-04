//go:build gui && windows

package main

import (
	"io"
	"log"
	"os"
	"syscall"
	"unsafe"
)

// Impede segunda instância (evita duas janelas GUI e, com PE errado, dois consoles).
const singleInstanceMutexName = `Local\VECX-VecxAgent-v1`

// ERROR_ALREADY_EXISTS (winerror.h); evita import de golang.org/x/sys/windows (melhor para gopls em Linux).
const errMutexAlreadyExists = syscall.Errno(183)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procFreeConsole  = kernel32.NewProc("FreeConsole")
	procCreateMutexW = kernel32.NewProc("CreateMutexW")
)

func init() {
	if procFreeConsole.Find() == nil {
		_, _, _ = procFreeConsole.Call()
	}
	log.SetOutput(io.Discard)

	name, err := syscall.UTF16PtrFromString(singleInstanceMutexName)
	if err != nil {
		return
	}
	r0, _, errno := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(name)))
	h := syscall.Handle(r0)
	if h == 0 {
		return
	}
	if errno == errMutexAlreadyExists {
		_ = syscall.CloseHandle(h)
		os.Exit(0)
	}
}
