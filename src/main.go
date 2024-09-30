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

const SNAKE_INITIAL_LENGTH = 5
const SNAKE_MAX_LENGTH = 256
const SNAKE_EYE_SIZE = 4
const SNAKE_EYE_OFFSET = 6

const GRID_CELL_SIZE = 32
const GRID_WIDTH = WINDOW_WIDTH / GRID_CELL_SIZE
const GRID_HEIGHT = WINDOW_HEIGHT / GRID_CELL_SIZE

const UPDATE_FRAME_RATE = 10
const SNAKE_OFFSET_DELTA = float32(GRID_CELL_SIZE) / float32(UPDATE_FRAME_RATE)

const FOOD_MAX_LENGTH = 5

type SnakeData struct {
	Position Point
	Velocity Point
	Offset   Vec2
}

type FoodData struct {
	Position Point
}

type GameState struct {
	FrameCount        uint64
	Snake             []SnakeData
	SnakeLength       uint8
	SnakeNextVelocity Point
	GrowSnake         bool
	Food              []FoodData
	FoodLength        uint8
}

func isSnakeAtPosition(gameState *GameState, position Point) bool {
	for i := 0; i < int(gameState.SnakeLength); i += 1 {
		if gameState.Snake[i].Position.Equals(position) {
			return true
		}
	}

	return false
}

func getFoodByPosition(gameState *GameState, position Point) *FoodData {
	for i := 0; i < int(gameState.FoodLength); i += 1 {
		if gameState.Food[i].Position.Equals(position) {
			return &gameState.Food[i]
		}
	}

	return nil
}

func getRandomPointNotContainingSnake(gameState *GameState) Point {
	point := RandomPoint(GRID_WIDTH, GRID_HEIGHT)
	for isSnakeAtPosition(gameState, point) {
		point = RandomPoint(GRID_WIDTH, GRID_HEIGHT)
	}
	return point
}

func main() {
	os.Exit(run())
}

func run() int {
	err := SDL.Init(SDL.INIT_VIDEO)

	if err != nil {
		SDL.LogError(SDL.LOG_CATEGORY_APPLICATION, "Failed initialize SDL: %s", err)
		return 1
	}
	defer SDL.Quit()

	SDL.LogSetPriority(SDL.LOG_CATEGORY_APPLICATION, SDL.LOG_PRIORITY_DEBUG)

	window, renderer, err := SDL.CreateWindowAndRenderer(WINDOW_WIDTH, WINDOW_HEIGHT, SDL.WINDOW_SHOWN)

	if err != nil {
		SDL.LogError(SDL.LOG_CATEGORY_APPLICATION, "Failed to create window and renderer: %s", err)
		return 1
	}

	defer SDL.DestroyWindow(window)
	defer SDL.DestroyRenderer(renderer)

	SDL.SetWindowTitle(window, WINDOW_TITLE)

	gameState := initialize()

	shouldQuit := false

	currentTime := SDL.GetTicks64()
	lastTime := currentTime

	for !shouldQuit {
		for event := SDL.PollEvent(); event != nil; event = SDL.PollEvent() {
			switch event.Type() {

			case SDL.QUIT:
				shouldQuit = true

			case SDL.KEYDOWN:
				switch event.Key().Keysym().Scancode {
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
		SnakeLength: SNAKE_INITIAL_LENGTH,
		Food:        make([]FoodData, FOOD_MAX_LENGTH),
	}

	for i := int32(0); i < int32(gameState.SnakeLength); i += 1 {
		gameState.Snake[i].Position.X = int32(gameState.SnakeLength) - 1 - i
		gameState.Snake[i].Velocity = Point{1, 0}
	}

	gameState.SnakeNextVelocity.X = 1

	foodPosition := getRandomPointNotContainingSnake(&gameState)

	gameState.Food[0].Position = foodPosition
	gameState.FoodLength = 1

	return gameState
}

func update(elapsedTime float32, gameState *GameState) {
	gameState.FrameCount += 1

	if gameState.FrameCount%UPDATE_FRAME_RATE == 0 {
		updateSnakePosition(gameState)
		checkForCollisions(gameState)
	} else {
		updateSnakeOffset(gameState)
	}
}

func updateSnakePosition(gameState *GameState) {
	if gameState.GrowSnake {
		gameState.SnakeLength += 1
		gameState.GrowSnake = false
	}

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
}

func checkForCollisions(gameState *GameState) {
	food := getFoodByPosition(gameState, gameState.Snake[0].Position)

	if food != nil {
		food.Position = getRandomPointNotContainingSnake(gameState)
		gameState.GrowSnake = true
	}
}

func updateSnakeOffset(gameState *GameState) {
	for i := 0; i < int(gameState.SnakeLength); i += 1 {
		gameState.Snake[i].Offset.X += SNAKE_OFFSET_DELTA * float32(gameState.Snake[i].Velocity.X)
		gameState.Snake[i].Offset.Y += SNAKE_OFFSET_DELTA * float32(gameState.Snake[i].Velocity.Y)
	}
}

func draw(renderer *SDL.Renderer, gameState *GameState) {
	rect := SDL.Rect{}

	SDL.SetRenderDrawColor(renderer, 100, 149, 237, 255)
	SDL.RenderClear(renderer)

	SDL.SetRenderDrawColor(renderer, 80, 80, 80, 255)
	for x := 0; x < WINDOW_WIDTH; x += GRID_CELL_SIZE {
		SDL.RenderDrawLine(renderer, x, 0, x, WINDOW_HEIGHT)
	}
	for y := 0; y < WINDOW_HEIGHT; y += GRID_CELL_SIZE {
		SDL.RenderDrawLine(renderer, 0, y, WINDOW_WIDTH, y)
	}

	// Food
	for i := 0; i < int(gameState.FoodLength); i += 1 {
		rect := SDL.Rect{
			X: gameState.Food[i].Position.X * GRID_CELL_SIZE,
			Y: gameState.Food[i].Position.Y * GRID_CELL_SIZE,
			W: GRID_CELL_SIZE,
			H: GRID_CELL_SIZE,
		}

		SDL.SetRenderDrawColor(renderer, 255, 0, 0, 255)
		SDL.RenderFillRect(renderer, &rect)
	}

	// Snake
	for i := (int)(gameState.SnakeLength - 1); i >= 0; i -= 1 {
		rect = SDL.Rect{
			X: gameState.Snake[i].Position.X*GRID_CELL_SIZE + int32(gameState.Snake[i].Offset.X),
			Y: gameState.Snake[i].Position.Y*GRID_CELL_SIZE + int32(gameState.Snake[i].Offset.Y),
			W: GRID_CELL_SIZE,
			H: GRID_CELL_SIZE,
		}

		color := (uint8)(float32(i+1) / float32(gameState.SnakeLength) * float32(255))
		SDL.SetRenderDrawColor(renderer, color, color, color, 255)
		SDL.RenderFillRect(renderer, &rect)
	}

	// Snake Eyes
	rect = SDL.Rect{
		X: gameState.Snake[0].Position.X*GRID_CELL_SIZE + int32(gameState.Snake[0].Offset.X) + SNAKE_EYE_OFFSET,
		Y: gameState.Snake[0].Position.Y*GRID_CELL_SIZE + int32(gameState.Snake[0].Offset.Y) + SNAKE_EYE_OFFSET,
		W: SNAKE_EYE_SIZE,
		H: SNAKE_EYE_SIZE,
	}

	if gameState.Snake[0].Velocity.X == RIGHT {
		rect.X += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	}

	if gameState.Snake[0].Velocity.Y == DOWN {
		rect.Y += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	}

	SDL.SetRenderDrawColor(renderer, 255, 255, 0, 255)
	SDL.RenderFillRect(renderer, &rect)

	if gameState.Snake[0].Velocity.X != 0 {
		rect.Y += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	}

	if gameState.Snake[0].Velocity.Y != 0 {
		rect.X += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	}

	SDL.RenderFillRect(renderer, &rect)

	SDL.RenderPresent(renderer)
}

func shutdown(gameState *GameState) {
}
