package hook

type Focus struct {
	focused *Scope
}

func (f *Focus) SetFocus(s *Scope) {
	f.focused = s
}

func (f *Focus) IsFocused(s *Scope) bool {
	return f.focused == s
}

func (f *Focus) IsInChain(s *Scope) bool {
	for cursor := f.focused; cursor != nil; cursor = cursor.parent {
		if cursor == s {
			return true
		}
	}

	return false
}
