package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("day22/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	depth, target := parseInput(string(data))

	cave := makeCaveSystem(depth, target)
	riskLevel := calculateRiskLevel(cave, target)
	log.Println("Total Risk Level:", riskLevel)
}

func parseInput(data string) (int, [2]int) {
	lines := strings.Split(data, "\n")
	depth, _ := strconv.Atoi(strings.Split(lines[0], " ")[1])
	coords := strings.Split(lines[1], " ")[1]
	parts := strings.Split(coords, ",")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	return depth, [2]int{x, y}
}

func makeCaveSystem(depth int, target [2]int) [][]int {
	cave := make([][]int, target[1]+1)
	for y := range cave {
		cave[y] = make([]int, target[0]+1)
		for x := range cave[y] {
			var geologicIndex int
			if x == 0 && y == 0 || x == target[0] && y == target[1] {
				geologicIndex = 0
			} else if y == 0 {
				geologicIndex = x * 16807
			} else if x == 0 {
				geologicIndex = y * 48271
			} else {
				geologicIndex = cave[y][x-1] * cave[y-1][x]
			}
			cave[y][x] = (geologicIndex + depth) % 20183
		}
	}
	return cave
}

func calculateRiskLevel(cave [][]int, target [2]int) int {
	riskLevel := 0
	for y := 0; y <= target[1]; y++ {
		for x := 0; x <= target[0]; x++ {
			riskLevel += cave[y][x] % 3
		}
	}
	return riskLevel
}
