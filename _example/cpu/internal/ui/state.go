package ui

type Focus int

const (
	FocusSearch Focus = iota
	FocusTable
	focusSentinel
	FocusPopup
)

type ViewState struct {
	CurrentFocus Focus
	Search       Search
	Searching    bool
	Table        Table
	Popup        Popup
}

func New() ViewState {
	return ViewState{
		CurrentFocus: FocusSearch,
		Search:       NewSearch(),
		Table:        NewTable(),
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
	s.Searching = searching
	s.Search.SetSearching(searching)
}
