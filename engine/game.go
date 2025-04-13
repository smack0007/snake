package engine

import (
	graphics "github.com/smack0007/snake/engine/graphics"
)

type Game interface {
	Initialize()

	Update(elapsedTime float32)

	Draw(graphics *graphics.Graphics)

	Shutdown()
}

func Run(game Game) error {
	platform, err := platformInitialize()

	if err != nil {
		return err
	}

	game.Initialize()

	platformRun(platform, game)

	game.Shutdown()

	platformShutdown(platform)

	return nil
}
