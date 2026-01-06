package mvu

import (
	"os"
	"os/exec"
)

type Event any

type shutdownEvent struct{}

var ShutdownEvent = shutdownEvent{}

type LaunchEvent struct {
	CmdBuilder func(ttyIn, ttyOut *os.File) (cmd *exec.Cmd, captureOutput bool, err error)
	OnResult   func(out []byte, err error) Task
}

type ResizeEvent struct {
	Width, Height int
}

type BatchTaskEvent struct {
	Tasks []Task
}
