package process

import (
	"sync"
)

type Profile struct {
	Info    *Info
	History *History
}

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

		if info.Stats != nil {
			s.profiles[info.PID].History.AddSample(*info.Stats)
			s.profiles[info.PID].Info.Stats = info.Stats
		}
	}

	for pid, _ := range s.profiles {
		if _, found := seenPIDs[pid]; !found {
			delete(s.profiles, pid)
		}
	}
}

func (s *Store) GetAll() []Profile {
	s.mu.RLock()
	defer s.mu.RUnlock()

	profiles := make([]Profile, 0, len(s.profiles))

	for _, profile := range s.profiles {
		info := *profile.Info
		history := &History{
			MaxSamples: profile.History.MaxSamples,
			Samples:    make([]Sample, len(profile.History.Samples)),
		}

		copy(history.Samples, profile.History.Samples)

		profiles = append(profiles, Profile{
			Info:    &info,
			History: history,
		})
	}

	return profiles
}
