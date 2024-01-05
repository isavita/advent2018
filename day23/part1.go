package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Nanobot struct {
	X, Y, Z, Radius int
}

func main() {
	file, err := os.Open("day23/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	nanobots := parseNanobots(file)

	strongest := findStrongestNanobot(nanobots)
	inRangeCount := countNanobotsInRange(nanobots, strongest)

	fmt.Println(inRangeCount)
}

func parseNanobots(file *os.File) []Nanobot {
	scanner := bufio.NewScanner(file)
	nanobots := []Nanobot{}
	re := regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		z, _ := strconv.Atoi(matches[3])
		radius, _ := strconv.Atoi(matches[4])

		nanobots = append(nanobots, Nanobot{x, y, z, radius})
	}

	return nanobots
}

func findStrongestNanobot(nanobots []Nanobot) Nanobot {
	var strongest Nanobot
	for _, nanobot := range nanobots {
		if nanobot.Radius > strongest.Radius {
			strongest = nanobot
		}
	}
	return strongest
}

func countNanobotsInRange(nanobots []Nanobot, strongest Nanobot) int {
	count := 0
	for _, nanobot := range nanobots {
		if manhattanDistance(nanobot, strongest) <= strongest.Radius {
			count++
		}
	}
	return count
}

func manhattanDistance(a, b Nanobot) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
