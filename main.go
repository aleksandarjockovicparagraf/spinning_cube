package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

type Vertex struct {
	X, Y, Z float64
}

var vertices = []Vertex{
	{0, 0, 0},
	{1, 0, 0},
	{1, 1, 0},
	{0, 1, 0},
	{0, 0, 1},
	{1, 0, 1},
	{1, 1, 1},
	{0, 1, 1},
}

var edges = [][]int{
	{0, 1}, {1, 2}, {2, 3}, {3, 0},
	{4, 5}, {5, 6}, {6, 7}, {7, 4},
	{0, 4}, {1, 5}, {2, 6}, {3, 7},
}

func main() {
	angleX := 0.0
	angleY := 0.0
	angleZ := 0.0
	theta := 0.001333

	center := calculateCenter(vertices)

	for {
		clearTerminal()

		for i := range vertices {
			vertices[i].RotateX(angleX, center)
			vertices[i].RotateY(angleY, center)
			vertices[i].RotateZ(angleZ, center)
		}

		drawCube(vertices)

		angleX += theta
		angleY += theta
		angleZ += theta

		time.Sleep(time.Millisecond * 70)
	}
}

func clearTerminal() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func calculateCenter(vertices []Vertex) Vertex {
	var center Vertex
	for _, v := range vertices {
		center.X += v.X
		center.Y += v.Y
		center.Z += v.Z
	}
	count := float64(len(vertices))
	center.X /= count
	center.Y /= count
	center.Z /= count
	return center
}

func (v *Vertex) RotateY(theta float64, center Vertex) {
	v.X -= center.X
	v.Z -= center.Z

	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)
	x := v.X*cosTheta + v.Z*sinTheta
	z := -v.X*sinTheta + v.Z*cosTheta
	v.X, v.Z = x, z

	v.X += center.X
	v.Z += center.Z
}

func (v *Vertex) RotateX(theta float64, center Vertex) {
	v.Y -= center.Y
	v.Z -= center.Z

	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)
	y := v.Y*cosTheta - v.Z*sinTheta
	z := v.Y*sinTheta + v.Z*cosTheta
	v.Y, v.Z = y, z

	v.Y += center.Y
	v.Z += center.Z
}

func (v *Vertex) RotateZ(theta float64, center Vertex) {
	v.X -= center.X
	v.Y -= center.Y

	cosTheta := math.Cos(theta)
	sinTheta := math.Sin(theta)
	x := v.X*cosTheta - v.Y*sinTheta
	y := v.X*sinTheta + v.Y*cosTheta
	v.X, v.Y = x, y

	v.X += center.X
	v.Y += center.Y
}

func drawCube(vertices []Vertex) {
	scale := 70
	projected := make([]Vertex, len(vertices))

	screenCenterX := float64(scale) / 2.0
	screenCenterY := float64(scale) / 2.0

	center := calculateCenter(vertices)

	offsetX := screenCenterX - center.X*float64(scale)/(center.Z+2)
	offsetY := screenCenterY - center.Y*float64(scale)/(center.Z+2)

	for i, v := range vertices {
		projected[i] = Vertex{
			X: (v.X*float64(scale)/(v.Z+2) + offsetX),
			Y: (v.Y*float64(scale)/(v.Z+2) + offsetY),
			Z: v.Z,
		}
	}

	grid := make([][]int, scale)
	for i := range grid {
		grid[i] = make([]int, scale)
	}

	for _, edge := range edges {
		v1 := projected[edge[0]]
		v2 := projected[edge[1]]
		drawLine(v1, v2, scale, grid)
	}

	printGrid(grid)
}

func drawLine(v1, v2 Vertex, scale int, grid [][]int) {
	startX := int(v1.X)
	startY := int(v1.Y)
	endX := int(v2.X)
	endY := int(v2.Y)

	dx := int(math.Abs(float64(endX - startX)))
	dy := int(math.Abs(float64(endY - startY)))
	sx, sy := 1, 1
	if startX > endX {
		sx = -1
	}
	if startY > endY {
		sy = -1
	}
	err := dx - dy

	x, y := startX, startY
	for {
		if x >= 0 && x < scale && y >= 0 && y < scale {
			grid[y][x] = 1
		}
		if x == endX && y == endY {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x += sx
		}
		if e2 < dx {
			err += dx
			y += sy
		}
	}
}

func printGrid(grid [][]int) {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 1 {
				fmt.Print("■")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
