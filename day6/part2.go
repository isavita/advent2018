package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	// Read input from file
	content, err := os.ReadFile("day6/input.txt")
	if err != nil {
		panic(err)
	}

	// Parse coordinates
	coordinates := parseCoordinates(string(content))

	// Find the region size
	regionSize := findRegionSize(coordinates, 10000)

	// Print the result
	fmt.Println(regionSize)
}

type coordinate struct {
	x, y int
}

func parseCoordinates(input string) []coordinate {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	coordinates := make([]coordinate, len(lines))

	for i, line := range lines {
		fmt.Sscanf(line, "%d, %d", &coordinates[i].x, &coordinates[i].y)
	}

	return coordinates
}

func findRegionSize(coordinates []coordinate, maxDistance int) int {
	minX, minY, maxX, maxY := findBoundingBox(coordinates)
	regionSize := 0

	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			totalDistance := 0

			for _, c := range coordinates {
				totalDistance += manhattanDistance(x, y, c.x, c.y)
			}

			if totalDistance < maxDistance {
				regionSize++
			}
		}
	}

	return regionSize
}

func findBoundingBox(coordinates []coordinate) (int, int, int, int) {
	minX, minY := math.MaxInt32, math.MaxInt32
	maxX, maxY := math.MinInt32, math.MinInt32

	for _, c := range coordinates {
		if c.x < minX {
			minX = c.x
		}
		if c.y < minY {
			minY = c.y
		}
		if c.x > maxX {
			maxX = c.x
		}
		if c.y > maxY {
			maxY = c.y
		}
	}

	return minX, minY, maxX, maxY
}

func manhattanDistance(x1, y1, x2, y2 int) int {
	return int(math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2)))
}
