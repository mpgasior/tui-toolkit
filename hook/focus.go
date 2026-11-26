package hook

type FocusManager struct {
	target *Scope
}

func (fm *FocusManager) SetFocus(s *Scope) {
	fm.target = s
}

func (fm *FocusManager) IsFocused(s *Scope) bool {
	return fm.target == s
}

func (fm *FocusManager) IsInChain(s *Scope) bool {
	for cursor := fm.target; cursor != nil; cursor = cursor.parent {
		if cursor == s {
			return true
		}
	}

	return false
}

func Focus(s *Scope) (focused bool, setFocus func()) {
	focused = s.focusMgr.IsFocused(s)
	setFocus = func() {
		s.focusMgr.SetFocus(s)
	}

	return focused, setFocus
}

type FocusChildFunc func(id string)

func FocusChild(s *Scope) FocusChildFunc {
	return func(id string) {
		if child, ok := s.children[id]; ok {
			s.focusMgr.SetFocus(child)
		}
	}
}
