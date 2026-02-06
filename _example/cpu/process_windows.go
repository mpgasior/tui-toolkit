package main

import (
	"sort"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func getProcessInfo(entry windows.ProcessEntry32) ProcessInfo {
	info := ProcessInfo{
		PID:  entry.ProcessID,
		Name: windows.UTF16ToString(entry.ExeFile[:]),
	}

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, entry.ProcessID)
	if err != nil {
		return info
	}
	defer windows.CloseHandle(handle)

	var creationTime, exitTime, kernelTime, userTime windows.Filetime

	if err := windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime); err != nil {
		return info
	}

	toDuration := func(ft windows.Filetime) time.Duration {
		ns := (int64(ft.HighDateTime) << 32) + int64(ft.LowDateTime)
		return time.Duration(ns * 100)
	}

	info.KernelTime = toDuration(kernelTime)
	info.UserTime = toDuration(userTime)

	return info
}

func ListProcesses() ([]ProcessInfo, error) {
	var processList []ProcessInfo

	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return nil, err
	}

	for {
		processList = append(processList, getProcessInfo(entry))

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}

	sort.Slice(processList, func(i, j int) bool {
		return processList[i].UserTime > processList[j].UserTime
	})

	return processList, nil
}
