package main

import (
	"math/rand/v2"
)

const (
	UP    = -1
	RIGHT = 1
	DOWN  = 1
	LEFT  = -1
)

type Point struct {
	X, Y int32
}

func (p1 Point) Equals(p2 Point) bool {
	return p1.X == p2.X && p1.Y == p2.Y
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

func (v1 Vec2) Equals(v2 Vec2) bool {
	return v1.X == v2.X && v1.Y == v2.Y
}
