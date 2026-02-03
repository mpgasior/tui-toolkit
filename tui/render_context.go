package tui

type RenderContext struct {
	Viewport Viewport
	Focused  bool
	Theme    Theme
}

func (rc RenderContext) WithViewport(vp Viewport) RenderContext {
	rc.Viewport = vp
	return rc
}

func (rc RenderContext) WithFocus(focus bool) RenderContext {
	rc.Focused = focus
	return rc
}
