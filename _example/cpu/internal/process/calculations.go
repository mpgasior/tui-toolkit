package process

import "runtime"

func CalculateCPU(first, last Sample) float64 {
	deltaWork := (last.UserTime + last.KernelTime) - (first.UserTime + first.KernelTime)
	deltaTime := last.SampleTime.Sub(first.SampleTime)
	if deltaTime <= 0 {
		return 0
	}

	rawUsage := float64(deltaWork) / float64(deltaTime)
	return rawUsage / float64(runtime.NumCPU()) * 100
}
