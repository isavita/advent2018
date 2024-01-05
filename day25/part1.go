package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z, t int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattanDistance(a, b Point) int {
	return abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z) + abs(a.t-b.t)
}

type UnionFind struct {
	parent []int
}

func NewUnionFind(size int) *UnionFind {
	parent := make([]int, size)
	for i := range parent {
		parent[i] = i
	}
	return &UnionFind{parent}
}

func (uf *UnionFind) find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.find(uf.parent[x])
	}
	return uf.parent[x]
}

func (uf *UnionFind) union(x, y int) {
	rootX := uf.find(x)
	rootY := uf.find(y)
	if rootX != rootY {
		uf.parent[rootX] = rootY
	}
}

func main() {
	file, err := os.Open("day25/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var points []Point
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		coords := strings.Split(scanner.Text(), ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		z, _ := strconv.Atoi(coords[2])
		t, _ := strconv.Atoi(coords[3])
		points = append(points, Point{x, y, z, t})
	}

	uf := NewUnionFind(len(points))
	for i := range points {
		for j := range points {
			if manhattanDistance(points[i], points[j]) <= 3 {
				uf.union(i, j)
			}
		}
	}

	constellationCount := 0
	for i, parent := range uf.parent {
		if i == parent {
			constellationCount++
		}
	}
	fmt.Println(constellationCount)
}
