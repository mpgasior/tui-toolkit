package ui

import (
	"github.com/mpgasior/tui-toolkit/_example/cpu/internal/model"
	"github.com/mpgasior/tui-toolkit/draw"
	"github.com/mpgasior/tui-toolkit/screen"
	"github.com/mpgasior/tui-toolkit/view"
)

type Popup struct {
	Result model.QueryResult
}

func (p *Popup) Draw(vp view.Port) {
	draw.Clear(vp, screen.DefaultStyle)
	draw.Box(vp, draw.BoxBorderDouble, screen.DefaultStyle.Fg(screen.ColorGreen))
	draw.Line(vp.Offset(0, 0, 0, 5), p.Result.Name, screen.DefaultStyle)
}
