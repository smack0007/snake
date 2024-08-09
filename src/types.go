package main

import (
	"math/rand/v2"
)

const (
	KEY_UP = iota
	KEY_RIGHT
	KEY_DOWN
	KEY_LEFT

	UP    = -1
	RIGHT = 1
	DOWN  = 1
	LEFT  = -1
)

type Point struct {
	X, Y int32
}

func (this Point) Equals(other Point) bool {
	return this.X == other.X && this.Y == other.Y
}

func RandomPoint(maxX int32, maxY int32) Point {
	return Point{
		X: rand.Int32() % maxX,
		Y: rand.Int32() % maxY,
	}
}

type Vec2 struct {
	X, Y float32
}

func (this Vec2) Equals(other Vec2) bool {
	return this.X == other.X && this.Y == other.Y
}
