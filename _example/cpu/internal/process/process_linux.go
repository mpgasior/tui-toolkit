package process

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetInfo(pid int) (Info, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return Info{}, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 15 {
		return Info{}, fmt.Errorf("insufficient data")
	}

	rawName := fields[1]
	name := strings.Trim(rawName, "()")

	utime, _ := strconv.ParseInt(fields[13], 10, 64)
	stime, _ := strconv.ParseInt(fields[14], 10, 64)

	const clockTicksPerSec = 100
	tickDuration := time.Second / clockTicksPerSec

	return Info{
		PID:  uint32(pid),
		Name: name,
		Stats: &Sample{
			UserTime:   time.Duration(utime) * tickDuration,
			KernelTime: time.Duration(stime) * tickDuration,
			SampleTime: time.Now().UTC(),
		},
	}, nil
}

func GetAll() ([]Info, error) {
	var processList []Info

	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		info, err := GetInfo(pid)
		if err == nil {
			processList = append(processList, info)
		}
	}

	return processList, nil
}
