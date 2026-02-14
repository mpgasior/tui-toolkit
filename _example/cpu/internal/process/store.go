package process

import (
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	profiles map[uint32]*Profile
}

func NewStore() *Store {
	return &Store{
		profiles: make(map[uint32]*Profile),
	}
}

func (s *Store) Sync(snapshot []Info) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seenPIDs := make(map[uint32]bool)
	for _, info := range snapshot {
		seenPIDs[info.PID] = true

		existing, found := s.profiles[info.PID]

		if !found || !existing.Info.CreationTime.Equal(info.CreationTime) {
			s.profiles[info.PID] = &Profile{
				Info:    &info,
				History: NewHistory(60),
			}
		}

		if info.LastSample != nil {
			s.profiles[info.PID].History.AddSample(*info.LastSample)
			s.profiles[info.PID].Info.LastSample = info.LastSample
		}
	}

	for pid, _ := range s.profiles {
		if _, found := seenPIDs[pid]; !found {
			delete(s.profiles, pid)
		}
	}
}

func (s *Store) GetAll() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots := make([]Snapshot, 0, len(s.profiles))

	for _, profile := range s.profiles {
		stats, computed := profile.History.Stats()
		snapshot := Snapshot{
			Info:      *profile.Info,
			AvgCPU:    stats.AvgCPU,
			RecentCPU: stats.RecentCPU,
			IsReady:   computed,
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}
