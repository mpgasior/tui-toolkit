package process

import (
	"time"
)

type Info struct {
	PID          uint32
	ParentPID    uint32
	Name         string
	CreationTime time.Time
	ExitTime     *time.Time
	Stats        *Sample
}
