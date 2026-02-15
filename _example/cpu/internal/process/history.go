package process

type Stats struct {
	RecentCPU float64
	AvgCPU    float64
}

type History struct {
	samples    []Sample
	maxSamples int
}

func NewHistory(maxSamples int) *History {
	return &History{
		samples:    make([]Sample, 0, maxSamples),
		maxSamples: maxSamples,
	}
}

func (h *History) Len() int {
	return len(h.samples)
}

func (h *History) Get(i int) Sample {
	return h.samples[i]
}

func (h *History) AddSample(s Sample) {
	if len(h.samples) >= h.maxSamples {
		copy(h.samples, h.samples[1:])
		h.samples[len(h.samples)-1] = s
		return
	}

	h.samples = append(h.samples, s)
}

func (h *History) Stats() (stats Stats, ok bool) {
	if h.Len() < 2 {
		return stats, false
	}
	stats = Stats{
		AvgCPU:    h.calculateAverage(),
		RecentCPU: h.calculateRecent(),
	}
	return stats, true
}

func (h *History) calculateAverage() float64 {
	if h.Len() < 2 {
		return 0.0
	}

	first := h.Get(0)
	last := h.Get(h.Len() - 1)

	return CalculateCPU(first, last)
}

func (h *History) calculateRecent() float64 {
	if h.Len() < 2 {
		return 0.0
	}

	first := h.Get(h.Len() - 2)
	last := h.Get(h.Len() - 1)

	return CalculateCPU(first, last)
}
