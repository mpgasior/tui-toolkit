package process

import (
	"runtime"
	"time"
)

type Sample struct {
	UserTime   time.Duration
	KernelTime time.Duration
	SampleTime time.Time
}

type Stats struct {
	AvgCPU    float64
	RecentCPU float64
}

type History struct {
	Samples    []Sample
	MaxSamples int
}

func NewHistory(maxSamples int) *History {
	return &History{
		Samples:    make([]Sample, 0, maxSamples),
		MaxSamples: maxSamples,
	}
}

func (h *History) Len() int {
	return len(h.Samples)
}

func (h *History) AddSample(s Sample) {
	if len(h.Samples) >= h.MaxSamples {
		copy(h.Samples, h.Samples[1:])
		h.Samples[len(h.Samples)-1] = s
		return
	}

	h.Samples = append(h.Samples, s)
}

func (h *History) Stats() Stats {
	return Stats{
		AvgCPU:    h.AvgCPU(),
		RecentCPU: h.RecentCPU(),
	}
}

func (h *History) AvgCPU() float64 {
	if len(h.Samples) < 2 {
		return 0
	}

	first := h.Samples[0]
	last := h.Samples[len(h.Samples)-1]

	deltaWork := (last.UserTime + last.KernelTime) - (first.UserTime + first.KernelTime)
	deltaTime := last.SampleTime.Sub(first.SampleTime)

	rawUsage := float64(deltaWork) / float64(deltaTime)
	return rawUsage / float64(runtime.NumCPU())
}

func (h *History) RecentCPU() float64 {
	if len(h.Samples) < 2 {
		return 0
	}

	first := h.Samples[len(h.Samples)-2]
	last := h.Samples[len(h.Samples)-1]

	deltaWork := (last.UserTime + last.KernelTime) - (first.UserTime + first.KernelTime)
	deltaTime := last.SampleTime.Sub(first.SampleTime)

	rawUsage := float64(deltaWork) / float64(deltaTime)
	return rawUsage / float64(runtime.NumCPU())
}
