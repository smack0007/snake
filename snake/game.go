package main

import (
	"math"

	graphics "github.com/smack0007/snake/engine/graphics"
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

type Game struct {
	FrameCount        uint64
	Snake             []SnakeData
	SnakeLength       uint8
	SnakeNextVelocity Point
	GrowSnake         bool
	Food              []FoodData
	FoodLength        uint8
}

func isSnakeAtPosition(game *Game, position Point) bool {
	for i := 0; i < int(game.SnakeLength); i += 1 {
		if game.Snake[i].Position.Equals(position) {
			return true
		}
	}

	return false
}

func getFoodByPosition(game *Game, position Point) *FoodData {
	for i := 0; i < int(game.FoodLength); i += 1 {
		if game.Food[i].Position.Equals(position) {
			return &game.Food[i]
		}
	}

	return nil
}

func getRandomPointNotContainingSnake(game *Game) Point {
	point := RandomPoint(GRID_WIDTH, GRID_HEIGHT)
	for isSnakeAtPosition(game, point) {
		point = RandomPoint(GRID_WIDTH, GRID_HEIGHT)
	}
	return point
}

func NewGame() Game {
	return Game{}
}

func (game *Game) Initialize() {
	game.FrameCount = 0
	game.Snake = make([]SnakeData, SNAKE_MAX_LENGTH)
	game.SnakeLength = SNAKE_INITIAL_LENGTH
	game.Food = make([]FoodData, FOOD_MAX_LENGTH)

	for i := int32(0); i < int32(game.SnakeLength); i += 1 {
		game.Snake[i].Position.X = int32(game.SnakeLength) - 1 - i
		game.Snake[i].Velocity = Point{1, 0}
	}

	game.SnakeNextVelocity.X = 1

	foodPosition := getRandomPointNotContainingSnake(game)

	game.Food[0].Position = foodPosition
	game.FoodLength = 1
}

func (game *Game) Update(elapsedTime float32) {
	game.FrameCount += 1

	if game.FrameCount%UPDATE_FRAME_RATE == 0 {
		updateSnakePosition(game)
		checkForCollisions(game)
	} else {
		updateSnakeOffset(game)
	}
}

func updateSnakePosition(game *Game) {
	if game.GrowSnake {
		game.SnakeLength += 1
		game.GrowSnake = false
	}

	for i := int(game.SnakeLength - 1); i >= 0; i -= 1 {
		if i == 0 {
			game.Snake[0].Position = Point{
				X: game.Snake[0].Position.X + game.Snake[0].Velocity.X,
				Y: game.Snake[0].Position.Y + game.Snake[0].Velocity.Y,
			}
		} else {
			game.Snake[i].Position = game.Snake[i-1].Position
			game.Snake[i].Velocity = game.Snake[i-1].Velocity
		}
		game.Snake[i].Offset = Vec2{0, 0}
	}

	game.Snake[0].Velocity = game.SnakeNextVelocity
}

func checkForCollisions(game *Game) {
	food := getFoodByPosition(game, game.Snake[0].Position)

	if food != nil {
		food.Position = getRandomPointNotContainingSnake(game)
		game.GrowSnake = true
	}
}

func updateSnakeOffset(game *Game) {
	for i := 0; i < int(game.SnakeLength); i += 1 {
		game.Snake[i].Offset.X += SNAKE_OFFSET_DELTA * float32(game.Snake[i].Velocity.X)
		game.Snake[i].Offset.Y += SNAKE_OFFSET_DELTA * float32(game.Snake[i].Velocity.Y)
	}
}

func (game *Game) Draw(g *graphics.Graphics) {
	// rect := SDL.Rect{}

	g.Clear(graphics.Color{
		R: 100,
		G: 149,
		B: 237,
		A: 255,
	})

	// SDL.SetRenderDrawColor(renderer, 80, 80, 80, 255)
	// for x := 0; x < WINDOW_WIDTH; x += GRID_CELL_SIZE {
	// 	SDL.RenderDrawLine(renderer, x, 0, x, WINDOW_HEIGHT)
	// }
	// for y := 0; y < WINDOW_HEIGHT; y += GRID_CELL_SIZE {
	// 	SDL.RenderDrawLine(renderer, 0, y, WINDOW_WIDTH, y)
	// }

	// // Food
	// for i := 0; i < int(gameState.FoodLength); i += 1 {
	// 	rect := SDL.Rect{
	// 		X: gameState.Food[i].Position.X * GRID_CELL_SIZE,
	// 		Y: gameState.Food[i].Position.Y * GRID_CELL_SIZE,
	// 		W: GRID_CELL_SIZE,
	// 		H: GRID_CELL_SIZE,
	// 	}

	// 	SDL.SetRenderDrawColor(renderer, 255, 0, 0, 255)
	// 	SDL.RenderFillRect(renderer, &rect)
	// }

	// // Snake
	// for i := (int)(gameState.SnakeLength - 1); i >= 0; i -= 1 {
	// 	rect = SDL.Rect{
	// 		X: gameState.Snake[i].Position.X*GRID_CELL_SIZE + int32(gameState.Snake[i].Offset.X),
	// 		Y: gameState.Snake[i].Position.Y*GRID_CELL_SIZE + int32(gameState.Snake[i].Offset.Y),
	// 		W: GRID_CELL_SIZE,
	// 		H: GRID_CELL_SIZE,
	// 	}

	// 	color := (uint8)(float32(i+1) / float32(gameState.SnakeLength) * float32(255))
	// 	SDL.SetRenderDrawColor(renderer, color, color, color, 255)
	// 	SDL.RenderFillRect(renderer, &rect)
	// }

	// // Snake Eyes
	// rect = SDL.Rect{
	// 	X: gameState.Snake[0].Position.X*GRID_CELL_SIZE + int32(gameState.Snake[0].Offset.X) + SNAKE_EYE_OFFSET,
	// 	Y: gameState.Snake[0].Position.Y*GRID_CELL_SIZE + int32(gameState.Snake[0].Offset.Y) + SNAKE_EYE_OFFSET,
	// 	W: SNAKE_EYE_SIZE,
	// 	H: SNAKE_EYE_SIZE,
	// }

	// if gameState.Snake[0].Velocity.X == RIGHT {
	// 	rect.X += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	// }

	// if gameState.Snake[0].Velocity.Y == DOWN {
	// 	rect.Y += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	// }

	// SDL.SetRenderDrawColor(renderer, 255, 255, 0, 255)
	// SDL.RenderFillRect(renderer, &rect)

	// if gameState.Snake[0].Velocity.X != 0 {
	// 	rect.Y += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	// }

	// if gameState.Snake[0].Velocity.Y != 0 {
	// 	rect.X += GRID_CELL_SIZE - SNAKE_EYE_SIZE - (SNAKE_EYE_OFFSET * 2)
	// }

	// SDL.RenderFillRect(renderer, &rect)

	// SDL.RenderPresent(renderer)
}

func (game *Game) Shutdown() {

}
