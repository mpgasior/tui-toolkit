package process

import (
	"sync"
	"time"
)

type Snapshot struct {
	Key Key

	PID          uint32
	Name         string
	CreationTime time.Time
	ExitTime     time.Time

	CPU       float64
	MemoryRSS uint64

	CPUAvg    float64
	MemoryMax uint64

	CPUReady bool
	MemReady bool
}

type HistorySnapshot struct {
	Key Key

	PID          uint32
	Name         string
	CreationTime time.Time
	ExitTime     time.Time

	CPU []float64
	Mem []uint64

	AvgCPU float64
	MaxMem uint64

	LatestCPU float64
	LatestMem uint64

	CPUReady bool
	MemReady bool
}

type Registry struct {
	mu sync.RWMutex

	table       map[Key]*Info
	telemetry   map[Key]*History
	historySize int
}

func NewRegistry(historySize int) *Registry {
	return &Registry{
		table:       make(map[Key]*Info),
		telemetry:   make(map[Key]*History),
		historySize: historySize,
	}
}

func (r *Registry) Update(samples []Sample) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	seenKeys := make(map[Key]bool, len(samples))

	for _, s := range samples {
		key := s.Key()
		seenKeys[key] = true

		if _, found := r.table[key]; !found {
			r.table[key] = &Info{
				PID:          s.PID,
				Name:         s.Name,
				CreationTime: s.CreationTime,
			}

			r.telemetry[key] = NewHistory(r.historySize)
		}

		r.telemetry[key].Push(s)
	}

	for k, v := range r.table {
		if _, ok := seenKeys[k]; ok {
			continue
		}

		if v.ExitTime.IsZero() {
			v.ExitTime = now
		}

		delete(r.table, k)
		delete(r.telemetry, k)
	}
}

func (r *Registry) GetSnapshot() []Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]Snapshot, 0, len(r.table))

	for k, i := range r.table {
		s := r.telemetry[k].Summary

		out = append(out, Snapshot{
			Key:          k,
			PID:          i.PID,
			Name:         i.Name,
			CreationTime: i.CreationTime,
			ExitTime:     i.ExitTime,
			CPU:          s.LatestCPU,
			CPUAvg:       s.AvgCPU,
			CPUReady:     s.CPUReady(),
			MemoryMax:    s.MaxMem,
			MemoryRSS:    s.LatestMem,
			MemReady:     s.MemReady(),
		})
	}

	return out
}

func (r *Registry) GetHistory(k Key) (HistorySnapshot, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, found := r.table[k]
	if !found {
		var zero HistorySnapshot
		return zero, false
	}

	h := r.telemetry[k]

	return HistorySnapshot{
		Key:          k,
		PID:          info.PID,
		Name:         info.Name,
		CreationTime: info.CreationTime,
		ExitTime:     info.ExitTime,

		CPU: h.CPU.All(),
		Mem: h.Mem.All(),

		LatestCPU: h.Summary.LatestCPU,
		LatestMem: h.Summary.LatestMem,
		AvgCPU:    h.Summary.AvgCPU,
		MaxMem:    h.Summary.MaxMem,
		CPUReady:  h.Summary.CPUReady(),
		MemReady:  h.Summary.MemReady(),
	}, true
}
