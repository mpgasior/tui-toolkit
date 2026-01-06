//go:build windows

package windowsx

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32          = windows.NewLazyDLL("kernel32.dll")
	procPeekConsoleIn = kernel32.NewProc("PeekConsoleInputW")
	procReadConsoleIn = kernel32.NewProc("ReadConsoleInputW")
)

func PeekConsoleInput(console windows.Handle, buffer []INPUT_RECORD) (n uint32, err error) {
	var numEventsRead uint32

	ok, _, err := procPeekConsoleIn.Call(
		uintptr(console),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&numEventsRead)),
	)

	if ok == 0 {
		return 0, err
	}

	return numEventsRead, nil
}

func ReadConsoleInput(console windows.Handle, buffer []INPUT_RECORD) (n uint32, err error) {
	var numEventsRead uint32

	ok, _, err := procReadConsoleIn.Call(
		uintptr(console),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		uintptr(unsafe.Pointer(&numEventsRead)),
	)

	if ok == 0 {
		return 0, err
	}

	return numEventsRead, nil
}
