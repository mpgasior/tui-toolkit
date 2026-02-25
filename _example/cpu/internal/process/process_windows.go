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

func GetSample(entry windows.ProcessEntry32) (Sample, error) {
	sample := Sample{
		PID:       entry.ProcessID,
		Name:      windows.UTF16ToString(entry.ExeFile[:]),
		Timestamp: time.Now().UTC(),
	}

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, entry.ProcessID)
	if err != nil {
		sample.IsRestricted = true
		return sample, err
	}
	defer windows.CloseHandle(handle)

	var creationTime, exitTime, kernelTime, userTime windows.Filetime

	if err := windows.GetProcessTimes(handle, &creationTime, &exitTime, &kernelTime, &userTime); err != nil {
		sample.IsRestricted = true
		return sample, err
	}

	counters, err := getMemoryInfo(handle)
	if err != nil {
		sample.IsRestricted = true
		return sample, err
	}

	toDuration := func(ft windows.Filetime) time.Duration {
		return time.Duration(uint64(ft.HighDateTime)<<32|uint64(ft.LowDateTime)) * 100
	}

	sample.KernelTotalTime = toDuration(kernelTime)
	sample.UserTotalTime = toDuration(userTime)
	sample.MemoryRSS = uint64(counters.WorkingSetSize)
	sample.CreationTime = time.Unix(0, creationTime.Nanoseconds())

	return sample, nil
}

func GetAll() ([]Sample, error) {
	var samples []Sample

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
		sample, _ := GetSample(entry)
		samples = append(samples, sample)

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}

	return samples, nil
}
