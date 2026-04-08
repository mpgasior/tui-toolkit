package mvu

import (
	"os"
)

type Event any

type shutdownEvent struct{}

var ShutdownEvent = shutdownEvent{}

type ExecEvent func(ttyIn, ttyOut *os.File) Task

type ResizeEvent struct {
	Width, Height int
}

type BatchTaskEvent struct {
	Tasks []Task
}
