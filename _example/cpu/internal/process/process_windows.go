package process

import (
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ProcessMemoryCounters struct {
	CB                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
}

var (
	psapi                = syscall.NewLazyDLL("psapi.dll")
	getProcessMemoryInfo = psapi.NewProc("GetProcessMemoryInfo")
)

func getMemoryInfo(process windows.Handle) (ProcessMemoryCounters, error) {
	var counters ProcessMemoryCounters
	counters.CB = uint32(unsafe.Sizeof(counters))

	ok, _, err := getProcessMemoryInfo.Call(
		uintptr(process),
		uintptr(unsafe.Pointer(&counters)),
		uintptr(counters.CB),
	)

	if ok == 0 {
		return counters, err
	}

	return counters, nil
}

func GetInfo(entry windows.ProcessEntry32) (Info, error) {
	info := Info{
		PID:       entry.ProcessID,
		ParentPID: entry.ParentProcessID,
		Name:      windows.UTF16ToString(entry.ExeFile[:]),
	}

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, entry.ProcessID)
	if err != nil {
		return info, err
	}
	defer windows.CloseHandle(handle)

	var creationTime, exitTime, kernelTime, userTime windows.Filetime

	if err := windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime); err != nil {
		return info, err
	}

	counters, err := getMemoryInfo(handle)
	if err != nil {
		return info, err
	}

	toDuration := func(ft windows.Filetime) time.Duration {
		return time.Duration(uint64(ft.HighDateTime)<<32|uint64(ft.LowDateTime)) * 100
	}

	info.LastSample = &Sample{
		KernelTime:     toDuration(kernelTime),
		UserTime:       toDuration(userTime),
		SampleTime:     time.Now().UTC(),
		WorkingSet:     uint64(counters.WorkingSetSize),
		VirtualSize:    uint64(counters.PagefileUsage),
		PeakWorkingSet: uint64(counters.PeakWorkingSetSize),
	}

	info.CreationTime = time.Unix(0, creationTime.Nanoseconds())
	if exitTime.HighDateTime != 0 && exitTime.LowDateTime != 0 {
		exit := time.Unix(0, exitTime.Nanoseconds())
		info.ExitTime = exit
	}

	return info, nil
}

func GetAll() ([]Info, error) {
	var processList []Info

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
		info, _ := GetInfo(entry)
		processList = append(processList, info)

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}

	return processList, nil
}
