package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getProcessInfo(pid int) (ProcessInfo, error) {
	data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return ProcessInfo{}, err
	}

	fields := strings.Fields(string(data))
	if len(fields) < 15 {
		return ProcessInfo{}, fmt.Errorf("insufficient data")
	}

	rawName := fields[1]
	name := strings.Trim(rawName, "()")

	utime, _ := strconv.ParseInt(fields[13], 10, 64)
	stime, _ := strconv.ParseInt(fields[14], 10, 64)

	const clockTicksPerSec = 100
	tickDuration := time.Second / clockTicksPerSec

	return ProcessInfo{
		PID:        uint32(pid),
		Name:       name,
		UserTime:   time.Duration(utime) * tickDuration,
		KernelTime: time.Duration(stime) * tickDuration,
	}, nil
}

func ListProcesses() ([]ProcessInfo, error) {
	var processList []ProcessInfo

	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		pid, err := strconv.Atoi(file.Name())
		if err != nil {
			continue
		}

		info, err := getProcessInfo(pid)
		if err == nil {
			processList = append(processList, info)
		}
	}

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].UserTime > processList[j].UserTime
	})

	return processList, nil
}
