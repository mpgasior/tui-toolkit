package tui

type RenderContext struct {
	Viewport Viewport
	Focused  bool
}

func (rc RenderContext) Fragment(x, y, w, h int) RenderContext {
	return RenderContext{
		Viewport: rc.Viewport.Slice(x, y, w, h),
		Focused:  rc.Focused,
	}
}

func (rc RenderContext) Unfocus() RenderContext {
	rc.Focused = false
	return rc
}
