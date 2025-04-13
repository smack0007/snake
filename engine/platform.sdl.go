//go:build sdl

package engine

import (
	"errors"
	"math"
	"runtime"
	"strings"

	graphics "github.com/smack0007/snake/engine/graphics"
	SDL "github.com/smack0007/snake/engine/sdl"
)

const WINDOW_WIDTH = 1024
const WINDOW_HEIGHT = 768
const WINDOW_TITLE = "Snake"
const DESIRED_FPS = 60

var GAME_TICK_RATE = (uint64)(math.Floor(float64(1000) / float64(DESIRED_FPS)))

type Platform struct {
	window *SDL.Window

	graphics *graphics.Graphics
}

func platformInitialize() (*Platform, error) {
	runtime.LockOSThread()

	err := SDL.Init(SDL.INIT_VIDEO)

	if err != nil {
		errMsg := strings.Join([]string{"Failed initialize SDL:", err.Error()}, " ")
		SDL.LogError(SDL.LOG_CATEGORY_APPLICATION, errMsg)
		return nil, errors.New(errMsg)
	}

	SDL.LogSetPriority(SDL.LOG_CATEGORY_APPLICATION, SDL.LOG_PRIORITY_DEBUG)

	window, renderer, err := SDL.CreateWindowAndRenderer(800, 600, SDL.WINDOW_SHOWN)

	if err != nil {
		errMsg := strings.Join([]string{"Failed to create window and renderer:", err.Error()}, " ")
		SDL.LogError(SDL.LOG_CATEGORY_APPLICATION, errMsg)
		return nil, errors.New(errMsg)
	}

	return &Platform{
		window: window,
		graphics: &graphics.Graphics{
			Renderer: renderer,
		},
	}, nil
}

func platformRun(platform *Platform, game Game) {
	shouldQuit := false

	currentTime := SDL.GetTicks64()
	lastTime := currentTime

	for !shouldQuit {
		for event := SDL.PollEvent(); event != nil; event = SDL.PollEvent() {
			switch event.Type() {

			case SDL.QUIT:
				shouldQuit = true

			case SDL.KEYDOWN:
				// switch event.Key().Keysym().Scancode {
				// case SDL.SCANCODE_UP:
				// 	if gameState.Snake[0].Velocity.Y != 1 {
				// 		gameState.SnakeNextVelocity = Point{X: 0, Y: -1}
				// 	}

				// case SDL.SCANCODE_RIGHT:
				// 	if gameState.Snake[0].Velocity.X != -1 {
				// 		gameState.SnakeNextVelocity = Point{X: 1, Y: 0}
				// 	}

				// case SDL.SCANCODE_DOWN:
				// 	if gameState.Snake[0].Velocity.Y != -1 {
				// 		gameState.SnakeNextVelocity = Point{X: 0, Y: 1}
				// 	}

				// case SDL.SCANCODE_LEFT:
				// 	if gameState.Snake[0].Velocity.X != 1 {
				// 		gameState.SnakeNextVelocity = Point{X: -1, Y: 0}
				// 	}
				// }
			}
		}

		currentTime = SDL.GetTicks64()
		elapsedTime := currentTime - lastTime

		if elapsedTime >= GAME_TICK_RATE {
			game.Update(float32(elapsedTime) / float32(1000))

			game.Draw(platform.graphics)

			SDL.RenderPresent(platform.graphics.Renderer)

			lastTime = currentTime
		}

		SDL.Delay(1)
	}
}

func platformShutdown(platform *Platform) {
	SDL.DestroyRenderer(platform.graphics.Renderer)
	SDL.DestroyWindow(platform.window)
	SDL.Quit()
}
