package ui

import "github.com/mpgasior/tui-toolkit/components"

type Focus int

const (
	FocusSearch Focus = iota
	FocusTable
	focusSentinel
)

type ViewState struct {
	CurrentFocus Focus
	Spinner      *components.Spinner
	TextInput    *components.TextInput
}

func New() ViewState {
	return ViewState{
		Spinner:   components.NewSpinner("spinner"),
		TextInput: &components.TextInput{},
	}
}

func (s *ViewState) NextFocus() {
	s.CurrentFocus = (s.CurrentFocus + 1) % focusSentinel
}

func (s *ViewState) PrevFocus() {
	s.CurrentFocus = (s.CurrentFocus - 1 + focusSentinel) % focusSentinel
}
