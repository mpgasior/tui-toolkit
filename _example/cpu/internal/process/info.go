package process

import (
	"fmt"
	"time"
)

type Key string

func (k Key) String() string {
	return string(k)
}

func NewKey(pid uint32, startTime time.Time) Key {
	return Key(fmt.Sprintf("%d-%d", pid, startTime.Unix()))
}

var KeyNone Key = ""

type Info struct {
	PID          uint32
	Name         string
	CreationTime time.Time
	ExitTime     time.Time
}

func (i Info) IsAlive() bool {
	return i.ExitTime.IsZero()
}

type Sample struct {
	PID             uint32
	Name            string
	CreationTime    time.Time
	UserTotalTime   time.Duration
	KernelTotalTime time.Duration
	MemoryRRS       uint64

	Timestamp    time.Time
	IsRestricted bool
}

func (s Sample) Key() Key {
	return NewKey(s.PID, s.CreationTime)
}
