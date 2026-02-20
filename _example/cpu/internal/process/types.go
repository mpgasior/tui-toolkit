package process

import (
	"time"
)

type Info struct {
	PID       uint32
	ParentPID uint32
	Name      string
}

type Profile struct {
	Info         Info
	History      History
	CreationTime time.Time
	Exited       bool
	ExitTime     time.Time
}

func (p *Profile) Clone() Profile {
	clone := *p
	clone.History = p.History.Clone()
	return clone
}

type Update struct {
	Info         Info
	Sample       *Sample
	CreationTime time.Time
	ExitTime     time.Time
}

type Sample struct {
	SampleTime time.Time

	UserTime       time.Duration
	KernelTime     time.Duration
	WorkingSet     uint64
	VirtualSize    uint64
	PeakWorkingSet uint64
}

type Snapshot struct {
	Info           Info
	CreationTime   time.Time
	ExitTime       time.Time
	Computed       bool
	AvgCPU         float64
	RecentCPU      float64
	WorkingSet     uint64
	PeakWorkingSet uint64
}
