package process

import (
	"sync"
)

type Store struct {
	mu           sync.RWMutex
	profiles     map[uint32]*Profile
	historyLimit int
}

func NewStore(historyLimit int) *Store {
	return &Store{
		profiles:     make(map[uint32]*Profile),
		historyLimit: historyLimit,
	}
}

func (s *Store) Sync(updates []Update) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seenPIDs := make(map[uint32]bool)
	for _, u := range updates {
		seenPIDs[u.Info.PID] = true

		existing, found := s.profiles[u.Info.PID]

		if !found || !existing.CreationTime.Equal(u.CreationTime) {
			existing = &Profile{
				Info:    u.Info,
				History: NewHistory(s.historyLimit),
			}
			s.profiles[u.Info.PID] = existing
		}

		if !u.CreationTime.IsZero() {
			existing.CreationTime = u.CreationTime
		}

		if !u.ExitTime.IsZero() {
			existing.ExitTime = u.ExitTime
		}

		if u.Sample != nil {
			existing.History.AddSample(*u.Sample)
		}
	}

	for pid, _ := range s.profiles {
		if _, found := seenPIDs[pid]; !found {
			delete(s.profiles, pid)
		}
	}
}

func (s *Store) GetProfile(pid uint32) (Profile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.profiles[pid]
	if !ok {
		return Profile{}, false
	}

	return p.Clone(), true
}

func (s *Store) GetAll() []Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshots := make([]Snapshot, 0, len(s.profiles))

	for _, profile := range s.profiles {
		stats, computed := profile.History.Stats()
		snapshot := Snapshot{
			Info:           profile.Info,
			CreationTime:   profile.CreationTime,
			ExitTime:       profile.ExitTime,
			AvgCPU:         stats.AvgCPU,
			RecentCPU:      stats.RecentCPU,
			Computed:       computed,
			WorkingSet:     stats.WorkingSet,
			PeakWorkingSet: stats.PeakWorkingSet,
		}

		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}
