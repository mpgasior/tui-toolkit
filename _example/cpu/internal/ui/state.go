package ui

type Focus int

const (
	FocusSearch Focus = iota
	FocusTable
	focusSentinel
)

type ViewState struct {
	CurrentFocus Focus
	Search       Search
	searching    bool
}

func New() ViewState {
	return ViewState{
		CurrentFocus: FocusSearch,
		Search:       NewSearch(),
	}
}

func (s *ViewState) IsFocused(f Focus) bool {
	return s.CurrentFocus == f
}

func (s *ViewState) NextFocus() {
	s.CurrentFocus = (s.CurrentFocus + 1) % focusSentinel
}

func (s *ViewState) PrevFocus() {
	s.CurrentFocus = (s.CurrentFocus - 1 + focusSentinel) % focusSentinel
}

func (s *ViewState) SetSearching(searching bool) {
	s.searching = searching
	s.Search.SetSearching(searching)
}
