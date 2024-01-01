package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	X, Y int
}

func main() {
	file, _ := os.Open("day6/input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	points := []Point{}
	maxX, maxY := 0, 0
	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ", ")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
		points = append(points, Point{x, y})
	}

	grid := make([][]int, maxX+2)
	for i := range grid {
		grid[i] = make([]int, maxY+2)
	}
	areas := make([]int, len(points))
	infinite := make([]bool, len(points))
	for i := range grid {
		for j := range grid[i] {
			minDist := maxX + maxY
			for k, point := range points {
				dist := abs(point.X-i) + abs(point.Y-j)
				if dist < minDist {
					minDist = dist
					grid[i][j] = k
				} else if dist == minDist {
					grid[i][j] = -1
				}
			}
			if grid[i][j] != -1 {
				if i == 0 || j == 0 || i == maxX+1 || j == maxY+1 {
					infinite[grid[i][j]] = true
				}
				areas[grid[i][j]]++
			}
		}
	}

	maxArea := 0
	for i, area := range areas {
		if !infinite[i] && area > maxArea {
			maxArea = area
		}
	}
	fmt.Println(maxArea)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
