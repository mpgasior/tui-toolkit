package mvu

import "github.com/mpgasior/tui-toolkit/view"

type RenderContext struct {
	View    view.Port
	Focused bool
}
