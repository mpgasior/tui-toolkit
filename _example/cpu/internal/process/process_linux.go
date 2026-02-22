package process

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Global or fetched once; typically 100 on most Linux systems
const clockTicksPerSec = 100

func parseStatusValue(valStr string) uint64 {
	fields := strings.Fields(valStr)
	if len(fields) == 0 {
		return 0
	}

	val, _ := strconv.ParseUint(fields[0], 10, 64)
	if len(fields) > 1 {
		switch fields[1] {
		case "kB":
			return val * 1024
		case "mB":
			return val * 1024 * 1024
		}
	}
	return val
}

// Helper to get system boot time to calculate absolute process start time
func getBootTime() time.Time {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return time.Unix(0, 0)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "btime") {
			f := strings.Fields(line)
			if len(f) > 1 {
				sec, _ := strconv.ParseInt(f[1], 10, 64)
				return time.Unix(sec, 0)
			}
		}
	}
	return time.Unix(0, 0)
}

func GetUpdate(pid int) (Update, error) {
	path := fmt.Sprintf("/proc/%d/status", pid)
	file, err := os.Open(path)
	if err != nil {
		return Update{}, err
	}
	defer file.Close()

	var (
		name      string
		ppid      uint32
		utime     int64
		stime     int64
		starttime int64
		currRSS   uint64
		peakRSS   uint64
	)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		val := strings.TrimSpace(parts[1])

		switch key {
		case "Name":
			name = val
		case "PPid":
			p, _ := strconv.Atoi(val)
			ppid = uint32(p)
		case "VmRSS":
			currRSS = parseStatusValue(val)
		case "VmHWM":
			peakRSS = parseStatusValue(val)
		}
	}

	statData, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err == nil {
		statFields := strings.Fields(string(statData))
		if len(statFields) > 21 {
			utime, _ = strconv.ParseInt(statFields[13], 10, 64)
			stime, _ = strconv.ParseInt(statFields[14], 10, 64)
			starttime, _ = strconv.ParseInt(statFields[21], 10, 64)
		}
	}

	tickDuration := time.Second / time.Duration(clockTicksPerSec)

	bootTime := getBootTime()
	creationTime := bootTime.Add(time.Duration(starttime) * tickDuration)

	return Update{
		Info: Info{
			PID:       uint32(pid),
			ParentPID: ppid,
			Name:      name,
		},
		CreationTime: creationTime.UTC(),
		Sample: &Sample{
			Timestamp:      time.Now().UTC(),
			UserTime:       time.Duration(utime) * tickDuration,
			KernelTime:     time.Duration(stime) * tickDuration,
			WorkingSet:     currRSS,
			PeakWorkingSet: peakRSS,
		},
	}, nil
}

func GetAll() ([]Update, error) {
	var updates []Update
	files, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if pid, err := strconv.Atoi(file.Name()); err == nil {
			if update, err := GetUpdate(pid); err == nil {
				updates = append(updates, update)
			}
		}
	}
	return updates, nil
}
