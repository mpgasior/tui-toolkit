package ui

import (
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/process"
)

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

func (s *ViewState) UpdateTable(rows []model.Process, query model.ListQuery) {
	s.Table.Set(rows, query.By, query.Order, query.Exclude)
}

func (s *ViewState) OpenPopup(key process.Key) {
	s.CurrentFocus = FocusPopup
	s.Popup.Open(key)
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
