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

func GetUpdate(entry windows.ProcessEntry32) (Update, error) {
	update := Update{
		Info: Info{
			PID:       entry.ProcessID,
			ParentPID: entry.ParentProcessID,
			Name:      windows.UTF16ToString(entry.ExeFile[:]),
		},
	}

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, entry.ProcessID)
	if err != nil {
		return update, err
	}
	defer windows.CloseHandle(handle)

	var creationTime, exitTime, kernelTime, userTime windows.Filetime

	if err := windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime); err != nil {
		return update, err
	}

	counters, err := getMemoryInfo(handle)
	if err != nil {
		return update, err
	}

	toDuration := func(ft windows.Filetime) time.Duration {
		return time.Duration(uint64(ft.HighDateTime)<<32|uint64(ft.LowDateTime)) * 100
	}

	update.Sample = &Sample{
		SampleTime:     time.Now().UTC(),
		KernelTime:     toDuration(kernelTime),
		UserTime:       toDuration(userTime),
		WorkingSet:     uint64(counters.WorkingSetSize),
		VirtualSize:    uint64(counters.PagefileUsage),
		PeakWorkingSet: uint64(counters.PeakWorkingSetSize),
	}

	update.CreationTime = time.Unix(0, creationTime.Nanoseconds())
	if exitTime.HighDateTime != 0 && exitTime.LowDateTime != 0 {
		exit := time.Unix(0, exitTime.Nanoseconds())
		update.ExitTime = exit
	}

	return update, nil
}

func GetAll() ([]Update, error) {
	var updates []Update

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
		update, _ := GetUpdate(entry)
		updates = append(updates, update)

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}

	return updates, nil
}
