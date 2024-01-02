package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("day11/input.txt")
	if err != nil {
		panic(err)
	}

	serial, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		panic(err)
	}

	const gridSize = 300
	grid := make([][]int, gridSize)
	for i := range grid {
		grid[i] = make([]int, gridSize)
	}

	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			rackID := x + 11
			powerLevel := rackID * (y + 1)
			powerLevel += serial
			powerLevel *= rackID
			powerLevel = (powerLevel / 100) % 10
			powerLevel -= 5
			grid[y][x] = powerLevel
		}
	}

	maxPower := -1 << 31
	var maxX, maxY int
	for y := 0; y < gridSize-2; y++ {
		for x := 0; x < gridSize-2; x++ {
			totalPower := 0
			for dy := 0; dy < 3; dy++ {
				for dx := 0; dx < 3; dx++ {
					totalPower += grid[y+dy][x+dx]
				}
			}
			if totalPower > maxPower {
				maxPower = totalPower
				maxX, maxY = x+1, y+1
			}
		}
	}

	fmt.Printf("%d,%d\n", maxX, maxY)
}
