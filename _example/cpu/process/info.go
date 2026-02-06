package process

import (
	"time"
)

type ProcessInfo struct {
	PID        uint32
	Name       string
	UserTime   time.Duration
	KernelTime time.Duration
}
