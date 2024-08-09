package main

type Point struct {
	X, Y int32
}

func (this Point) Equals(other Point) bool {
	return this.X == other.X && this.Y == other.Y
}

type Vec2 struct {
	X, Y float32
}

func (this Vec2) Equals(other Vec2) bool {
	return this.X == other.X && this.Y == other.Y
}
