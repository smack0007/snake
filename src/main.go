package main

import (
	"math"
	"os"

	SDL "github.com/smack0007/sdl-go/sdl"
)

const WINDOW_WIDTH = 1024
const WINDOW_HEIGHT = 768
const WINDOW_TITLE = "Snake"
const DESIRED_FPS = 60

var GAME_TICK_RATE = (uint64)(math.Floor(float64(1000) / float64(DESIRED_FPS)))

const SANKE_INITIAL_LENGTH = 5
const SNAKE_MAX_LENGTH = 256
const SNAKE_PART_SIZE = 32

const GRID_WIDTH = WINDOW_WIDTH / SNAKE_PART_SIZE
const GRID_HEIGHT = WINDOW_HEIGHT / SNAKE_PART_SIZE
const UPDATE_FRAME_RATE = 10
const SNAKE_OFFSET_DELTA = float32(SNAKE_PART_SIZE) / float32(UPDATE_FRAME_RATE)

const FOOD_MAX_LENGTH = 5

type SnakeData struct {
	Position Point
	Velocity Point
	Offset   Vec2
}

type FoodData struct {
}

type GameState struct {
	FrameCount        uint64
	Snake             []SnakeData
	SnakeLength       uint8
	SnakeNextVelocity Point
	Food              []FoodData
	FoodLength        uint8
}

func main() {
	os.Exit(run())
}

func run() int {
	if SDL.Init(SDL.INIT_VIDEO) != 0 {
		SDL.LogError(SDL.LOG_CATEGORY_APPLICATION, "Failed initialize SDL.")
		return 1
	}
	defer SDL.Quit()

	SDL.LogSetPriority(SDL.LOG_CATEGORY_APPLICATION, SDL.LOG_PRIORITY_DEBUG)

	var window *SDL.Window
	var renderer *SDL.Renderer
	result := SDL.CreateWindowAndRenderer(WINDOW_WIDTH, WINDOW_HEIGHT, SDL.WINDOW_SHOWN, &window, &renderer)

	if result != 0 {
		SDL.LogError(SDL.LOG_CATEGORY_APPLICATION, "Failed to create window and renderer.")
		return 1
	}

	defer SDL.DestroyWindow(window)
	defer SDL.DestroyRenderer(renderer)

	SDL.SetWindowTitle(window, WINDOW_TITLE)

	gameState := initialize()

	shouldQuit := false
	event := SDL.Event{}

	currentTime := SDL.GetTicks64()
	lastTime := currentTime

	for !shouldQuit {
		for SDL.PollEvent(&event) > 0 {
			switch event.Type() {

			case SDL.QUIT:
				shouldQuit = true

			case SDL.KEYDOWN:
				switch event.Key.Keysym().Scancode {
				case SDL.SCANCODE_UP:
					if gameState.Snake[0].Velocity.Y != 1 {
						gameState.SnakeNextVelocity = Point{X: 0, Y: -1}
					}

				case SDL.SCANCODE_RIGHT:
					if gameState.Snake[0].Velocity.X != -1 {
						gameState.SnakeNextVelocity = Point{X: 1, Y: 0}
					}

				case SDL.SCANCODE_DOWN:
					if gameState.Snake[0].Velocity.Y != -1 {
						gameState.SnakeNextVelocity = Point{X: 0, Y: 1}
					}

				case SDL.SCANCODE_LEFT:
					if gameState.Snake[0].Velocity.X != 1 {
						gameState.SnakeNextVelocity = Point{X: -1, Y: 0}
					}
				}
			}
		}

		currentTime = SDL.GetTicks64()
		elapsedTime := currentTime - lastTime

		if elapsedTime >= GAME_TICK_RATE {
			update(float32(elapsedTime)/float32(1000), &gameState)
			draw(renderer, &gameState)

			lastTime = currentTime
		}

		SDL.Delay(1)
	}

	shutdown(&gameState)

	return 0
}

func initialize() GameState {
	gameState := GameState{
		FrameCount:  0,
		Snake:       make([]SnakeData, SNAKE_MAX_LENGTH),
		SnakeLength: SANKE_INITIAL_LENGTH,
	}

	for i := int32(0); i < int32(gameState.SnakeLength); i += 1 {
		gameState.Snake[i].Position.X = int32(gameState.SnakeLength) - 1 - i
		gameState.Snake[i].Velocity = Point{1, 0}
	}

	gameState.SnakeNextVelocity.X = 1

	return gameState
}

func update(elapsedTime float32, gameState *GameState) {
	gameState.FrameCount += 1

	if gameState.FrameCount%UPDATE_FRAME_RATE == 0 {
		for i := int(gameState.SnakeLength - 1); i >= 0; i -= 1 {
			if i == 0 {
				gameState.Snake[0].Position = Point{
					X: gameState.Snake[0].Position.X + gameState.Snake[0].Velocity.X,
					Y: gameState.Snake[0].Position.Y + gameState.Snake[0].Velocity.Y,
				}
			} else {
				gameState.Snake[i].Position = gameState.Snake[i-1].Position
				gameState.Snake[i].Velocity = gameState.Snake[i-1].Velocity
			}
			gameState.Snake[i].Offset = Vec2{0, 0}
		}

		gameState.Snake[0].Velocity = gameState.SnakeNextVelocity
	} else {
		for i := 0; i < int(gameState.SnakeLength); i += 1 {
			gameState.Snake[i].Offset.X += SNAKE_OFFSET_DELTA * float32(gameState.Snake[i].Velocity.X)
			gameState.Snake[i].Offset.Y += SNAKE_OFFSET_DELTA * float32(gameState.Snake[i].Velocity.Y)
		}
	}
}

func draw(renderer *SDL.Renderer, gameState *GameState) {
	SDL.SetRenderDrawColor(renderer, 100, 149, 237, 255)
	SDL.RenderClear(renderer)

	for i := (int)(gameState.SnakeLength - 1); i >= 0; i -= 1 {
		color := (uint8)(float32(i+1) / float32(gameState.SnakeLength) * float32(255))
		SDL.SetRenderDrawColor(renderer, color, color, color, 255)
		rect := SDL.Rect{
			X: gameState.Snake[i].Position.X*SNAKE_PART_SIZE + int32(gameState.Snake[i].Offset.X),
			Y: gameState.Snake[i].Position.Y*SNAKE_PART_SIZE + int32(gameState.Snake[i].Offset.Y),
			W: SNAKE_PART_SIZE,
			H: SNAKE_PART_SIZE,
		}
		SDL.FillRect(renderer, &rect)
	}

	SDL.RenderPresent(renderer)
}

func shutdown(gameState *GameState) {
}
