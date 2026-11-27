package hook

type ScheduleFunc func()

func Root(schedule ScheduleFunc) *Scope {
	s := NewScope(nil)

	s.schedule = schedule
	s.focusMgr = &FocusManager{target: s}

	return s
}

func DispatchKey(s *Scope, e *KeyEvent) {
	for cursor := s.focusMgr.target; cursor != nil; cursor = cursor.parent {
		if cursor.keyHandlers != nil {
			for _, handler := range cursor.keyHandlers {
				handler(e)
			}

			if e.stopped {
				return
			}
		}
	}
}
