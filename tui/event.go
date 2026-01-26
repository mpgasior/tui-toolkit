package tui

type Event any

type BatchEvent struct {
	Events []Event
}

func Batch(events ...Event) BatchEvent {
	return BatchEvent{Events: events}
}
