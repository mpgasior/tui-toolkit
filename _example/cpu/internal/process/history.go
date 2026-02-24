package process

import (
	"runtime"
	"time"
)

type Summary struct {
	CPUSum float64

	LatestCPU float64
	AvgCPU    float64
	LatestMem uint64
	MaxMem    uint64

	Samples uint64
}

func (s Summary) CPUReady() bool {
	return true
	return s.Samples >= 2
}

func (s Summary) MemReady() bool {
	return true
	return s.Samples >= 1
}

type History struct {
	CPU *RingBuffer[float64]
	Mem *RingBuffer[uint64]

	lastTimestamp   time.Time
	lastUserTotal   time.Duration
	lastKernelTotal time.Duration

	Summary Summary
}

func NewHistory(size int) *History {
	return &History{
		CPU: NewRingBuffer[float64](size),
		Mem: NewRingBuffer[uint64](size),
	}
}

func (h *History) Push(s Sample) {
	h.Mem.Push(s.MemoryRRS)
	h.Summary.LatestMem = s.MemoryRRS
	if s.MemoryRRS > h.Summary.MaxMem {
		h.Summary.MaxMem = s.MemoryRRS
	}

	if h.lastTimestamp.IsZero() {
		h.lastTimestamp = s.Timestamp
		h.lastUserTotal = s.UserTotalTime
		h.lastKernelTotal = s.KernelTotalTime
		return
	}

	deltaTime := s.Timestamp.Sub(h.lastTimestamp)
	if deltaTime > 0 {
		lastTotal := h.lastKernelTotal + h.lastUserTotal
		total := s.KernelTotalTime + s.UserTotalTime

		deltaWork := total - lastTotal
		rawUsage := float64(deltaWork) / float64(deltaTime)
		usage := rawUsage / float64(runtime.NumCPU())
		usage *= 100

		if h.CPU.Full() {
			oldestUsage, _ := h.CPU.Get(0)
			h.Summary.CPUSum -= oldestUsage
		} else {
			h.Summary.Samples += 1
		}

		h.Summary.CPUSum += usage
		h.Summary.LatestCPU = usage
		h.CPU.Push(usage)

		h.Summary.AvgCPU = h.Summary.CPUSum / float64(h.Summary.Samples)
	}

	h.lastTimestamp = s.Timestamp
	h.lastUserTotal = s.UserTotalTime
	h.lastKernelTotal = s.KernelTotalTime
}
