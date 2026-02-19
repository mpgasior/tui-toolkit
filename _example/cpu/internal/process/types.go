package process

import (
	"time"
)

type Info struct {
	PID          uint32
	ParentPID    uint32
	Name         string
	CreationTime time.Time
	ExitTime     time.Time
	LastSample   *Sample
}

type Profile struct {
	Info    *Info
	History *History
}

type Sample struct {
	UserTime       time.Duration
	KernelTime     time.Duration
	SampleTime     time.Time
	WorkingSet     uint64
	VirtualSize    uint64
	PeakWorkingSet uint64
}

type Snapshot struct {
	Info           Info
	IsReady        bool
	AvgCPU         float64
	RecentCPU      float64
	WorkingSet     uint64
	PeakWorkingSet uint64
}
