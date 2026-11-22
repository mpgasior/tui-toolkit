package hook

import "slices"

type Focus struct {
	focused *Scope
}

func (f *Focus) SetFocus(s *Scope) {
	f.focused = s
}

func (f *Focus) Chain() []*Scope {
	var chain []*Scope
	for cursor := f.focused; cursor != nil; cursor = cursor.parent {
		if cursor.IsFocusable() {
			chain = append(chain, cursor)
		}
	}
	return chain
}

func (f *Focus) IsFocused(s *Scope) bool {
	return f.focused == s
}

func (f *Focus) IsInChain(s *Scope) bool {
	if s == nil {
		return false
	}

	return slices.Contains(f.Chain(), s)
}
