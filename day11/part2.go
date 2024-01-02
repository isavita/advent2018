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

	// Calculate power level for each cell
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
	var maxX, maxY, maxSize int
	for size := 1; size <= gridSize; size++ {
		for y := 0; y < gridSize-size+1; y++ {
			for x := 0; x < gridSize-size+1; x++ {
				totalPower := 0
				for dy := 0; dy < size; dy++ {
					for dx := 0; dx < size; dx++ {
						totalPower += grid[y+dy][x+dx]
					}
				}
				if totalPower > maxPower {
					maxPower = totalPower
					maxX, maxY, maxSize = x+1, y+1, size
				}
			}
		}
	}

	fmt.Printf("%d,%d,%d\n", maxX, maxY, maxSize)
}
