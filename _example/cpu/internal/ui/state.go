package ui

type Focus int

const (
	FocusSearch Focus = iota
	FocusTable
)

type ViewState struct {
	CurrentFocus Focus
}
