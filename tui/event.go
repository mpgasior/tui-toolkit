package tui

import "slices"

import "github.com/mpgasior/tui-go/vt"

type Event any

type shutdownEvent struct{}

var ShutdownEvent = shutdownEvent{}

type KeyEvent struct {
	Key  vt.Key
	Rune rune
}

func (e KeyEvent) IsKey(keys ...vt.Key) bool {
	return slices.Contains(keys, e.Key)
}

func (e KeyEvent) IsRune(r rune) bool {
	return e.Rune == r
}

type PasteEvent struct {
	Bytes []byte
}
