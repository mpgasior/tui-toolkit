package tui

import (
	"os"
	"os/exec"
	"slices"

	"github.com/mpgasior/tui-toolkit/vt"
)

type Event any

type shutdownEvent struct{}

var ShutdownEvent = shutdownEvent{}

type KeyEvent struct {
	Key  vt.Key
	Rune rune
}

func (e KeyEvent) IsKey(keys ...vt.Key) bool {
	return slices.Contains(keys, e.Key)
}

func (e KeyEvent) IsRune(r rune) bool {
	return e.Rune == r
}

type PasteEvent struct {
	Bytes []byte
}

type LaunchEvent struct {
	CmdBuilder func(ttyIn, ttyOut *os.File) (cmd *exec.Cmd, captureOutput bool, err error)
	OnResult   func(out []byte, err error) Task
}

type ResizeEvent struct {
	Width, Height int
}
