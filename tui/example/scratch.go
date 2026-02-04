package main

import (
	"fmt"
	"sort"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ProcessStats struct {
	Pid  uint32
	Name string
	CPU  float64
}

// snapshotCPU captures the total CPU time for all processes currently running
func snapshotCPU() (map[uint32]time.Duration, map[uint32]string) {
	times := make(map[uint32]time.Duration)
	names := make(map[uint32]string)

	// Take a snapshot of all processes
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return times, names
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	for err = windows.Process32First(snapshot, &entry); err == nil; err = windows.Process32Next(snapshot, &entry) {
		pid := entry.ProcessID
		names[pid] = windows.UTF16ToString(entry.ExeFile[:])

		// Open process to get timing info
		// PROCESS_QUERY_LIMITED_INFORMATION is less restrictive than PROCESS_QUERY_INFORMATION
		h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
		if err != nil {
			continue // Skip processes we can't touch (System/Protected)
		}

		var creation, exit, kernel, user windows.Filetime
		if err := windows.GetProcessTimes(h, &creation, &exit, &kernel, &user); err == nil {
			times[pid] = time.Duration(kernel.Nanoseconds() + user.Nanoseconds())
		}
		windows.CloseHandle(h)
	}

	return times, names
}

func main() {
	fmt.Println("Collecting initial snapshot...")
	t1, names := snapshotCPU()
	wall1 := time.Now()

	for {
		time.Sleep(1 * time.Second)

		t2, _ := snapshotCPU()
		wall2 := time.Now()
		deltaWall := wall2.Sub(wall1).Seconds()

		var results []ProcessStats
		for pid, totalTime2 := range t2 {
			if totalTime1, ok := t1[pid]; ok {
				deltaCPU := (totalTime2 - totalTime1).Seconds()
				usage := (deltaCPU / deltaWall) * 100

				if usage > 0.1 { // Only show processes doing actual work
					results = append(results, ProcessStats{pid, names[pid], usage})
				}
			}
		}

		// Sort by CPU usage descending
		sort.Slice(results, func(i, j int) bool {
			return results[i].CPU > results[j].CPU
		})

		// Print Results
		fmt.Print("\033[H\033[2J") // Clear terminal screen
		fmt.Printf("%-10s %-25s %-10s\n", "PID", "Name", "CPU %")
		fmt.Println("--------------------------------------------------")
		for _, p := range results {
			fmt.Printf("%-10d %-25s %6.2f%%\n", p.Pid, p.Name, p.CPU)
		}

		// Update snapshots for next loop
		t1, wall1 = t2, wall2
	}
}
