package mvu

type Component interface {
	Init() Task
	Update(e Event) Task
	Render(ctx RenderContext)
}
