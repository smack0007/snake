//go:build sdl

package graphics

import (
	SDL "github.com/smack0007/snake/engine/sdl"
)

type Graphics struct {
	Renderer *SDL.Renderer
}

func (graphics Graphics) Clear(color Color) {
	SDL.SetRenderDrawColor(graphics.Renderer, color.R, color.G, color.B, color.A)
	SDL.RenderClear(graphics.Renderer)
}
