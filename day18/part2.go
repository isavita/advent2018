package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Open       = '.'
	Trees      = '|'
	Lumberyard = '#'
	Size       = 50
)

func main() {
	grid, err := readInput("day18/input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	seenStates := make(map[string]int)
	var cycleStart, cycleLength int
	for minute := 0; ; minute++ {
		state := gridToString(grid)
		if seenMinute, found := seenStates[state]; found {
			cycleStart = seenMinute
			cycleLength = minute - seenMinute
			break
		}
		seenStates[state] = minute
		grid = transform(grid)
	}

	remainingMinutes := (1000000000 - cycleStart) % cycleLength
	for i := 0; i < remainingMinutes; i++ {
		grid = transform(grid)
	}

	wooded, lumberyards := countResources(grid)
	fmt.Println(wooded * lumberyards)
}

func readInput(filename string) ([][]rune, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var grid [][]rune
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]rune, len(line))
		for i, r := range line {
			row[i] = r
		}
		grid = append(grid, row)
	}
	return grid, scanner.Err()
}

func transform(grid [][]rune) [][]rune {
	newGrid := make([][]rune, len(grid))
	for i := range grid {
		newGrid[i] = make([]rune, len(grid[i]))
		for j := range grid[i] {
			newGrid[i][j] = nextAcreState(grid, i, j)
		}
	}
	return newGrid
}

func nextAcreState(grid [][]rune, i, j int) rune {
	switch grid[i][j] {
	case Open:
		if countAdjacent(grid, i, j, Trees) >= 3 {
			return Trees
		}
	case Trees:
		if countAdjacent(grid, i, j, Lumberyard) >= 3 {
			return Lumberyard
		}
	case Lumberyard:
		if countAdjacent(grid, i, j, Lumberyard) >= 1 && countAdjacent(grid, i, j, Trees) >= 1 {
			return Lumberyard
		}
		return Open
	}
	return grid[i][j]
}

func countAdjacent(grid [][]rune, i, j int, acreType rune) int {
	count := 0
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			if i+x >= 0 && i+x < len(grid) && j+y >= 0 && j+y < len(grid[i]) && grid[i+x][j+y] == acreType {
				count++
			}
		}
	}
	return count
}

func countResources(grid [][]rune) (int, int) {
	wooded, lumberyards := 0, 0
	for i := range grid {
		for j := range grid[i] {
			switch grid[i][j] {
			case Trees:
				wooded++
			case Lumberyard:
				lumberyards++
			}
		}
	}
	return wooded, lumberyards
}

func gridToString(grid [][]rune) string {
	var sb strings.Builder
	for _, row := range grid {
		sb.WriteString(string(row))
		sb.WriteRune('\n')
	}
	return sb.String()
}
