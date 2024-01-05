package main

import (
	"fmt"
	"log"
	"os"
)

type Point struct{ X, Y int }
type DoorMap map[Point]map[Point]bool

func main() {
	data, err := os.ReadFile("day20/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	regex := string(data)
	dm := buildMap(regex[1 : len(regex)-1])
	maxDoors := findFurthestRoom(dm)
	fmt.Println(maxDoors)
}

func buildMap(regex string) DoorMap {
	dm := make(DoorMap)
	var stack []Point
	cp := Point{0, 0}
	for _, c := range regex {
		if c == '(' {
			stack = append(stack, cp)
		} else if c == '|' {
			cp = stack[len(stack)-1]
		} else if c == ')' {
			cp = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		} else {
			np := move(cp, c)
			if dm[cp] == nil {
				dm[cp] = make(map[Point]bool)
			}
			dm[cp][np] = true
			cp = np
		}
	}
	return dm
}

func move(p Point, dir rune) Point {
	switch dir {
	case 'N':
		return Point{p.X, p.Y - 1}
	case 'S':
		return Point{p.X, p.Y + 1}
	case 'E':
		return Point{p.X + 1, p.Y}
	case 'W':
		return Point{p.X - 1, p.Y}
	}
	return p
}

func findFurthestRoom(dm DoorMap) int {
	visited := make(map[Point]int)
	var queue []Point
	queue = append(queue, Point{0, 0})
	maxDoors := 0

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for np := range dm[p] {
			if _, seen := visited[np]; !seen {
				visited[np] = visited[p] + 1
				maxDoors = max(maxDoors, visited[np])
				queue = append(queue, np)
			}
		}
	}
	return maxDoors
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
