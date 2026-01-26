package tui

type Component interface {
	Init() Task
	Update(e Event) Task
	Render(v View, focused bool)
}
