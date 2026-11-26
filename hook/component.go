package hook

import "github.com/mpgasior/tui-toolkit/draw"

type RenderFunc func(sf *draw.Surface)

type Component func(s *Scope) RenderFunc
