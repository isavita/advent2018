package main

import (
	"bufio"
	"container/heap"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Nanobot struct {
	X, Y, Z, R int
}

type Point struct {
	X, Y, Z int
}

type Region struct {
	Min, Max Point
	Count    int
}

type PriorityQueue []*Region

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Count == pq[j].Count {
		return manhattan(pq[i].Min, Point{}) < manhattan(pq[j].Min, Point{})
	}
	return pq[i].Count > pq[j].Count
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Region))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func main() {
	nanobots := readNanobots("day23/input.txt")
	result := findBestPoint(nanobots)
	println("Best point distance:", manhattan(result.Min, Point{}))
}

func readNanobots(filename string) []Nanobot {
	file, _ := os.Open(filename)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var nanobots []Nanobot
	re := regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(\d+)`)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		matches := re.FindStringSubmatch(text)
		x, _ := strconv.Atoi(matches[1])
		y, _ := strconv.Atoi(matches[2])
		z, _ := strconv.Atoi(matches[3])
		r, _ := strconv.Atoi(matches[4])
		nanobots = append(nanobots, Nanobot{x, y, z, r})
	}

	return nanobots
}

func findBestPoint(nanobots []Nanobot) *Region {
	var pq PriorityQueue
	heap.Init(&pq)

	// Initial large region to start the search
	initialRegion := &Region{
		Min:   Point{X: math.MinInt32, Y: math.MinInt32, Z: math.MinInt32},
		Max:   Point{X: math.MaxInt32, Y: math.MaxInt32, Z: math.MaxInt32},
		Count: len(nanobots),
	}
	heap.Push(&pq, initialRegion)

	var best *Region

	for pq.Len() > 0 {
		region := heap.Pop(&pq).(*Region)
		if best != nil && region.Count < best.Count {
			break
		}

		// Check if the region is a single point
		if region.Min == region.Max {
			if best == nil || region.Count > best.Count {
				best = region
			}
			continue
		}

		// Split the region into smaller subregions
		midX := (region.Min.X + region.Max.X) / 2
		midY := (region.Min.Y + region.Max.Y) / 2
		midZ := (region.Min.Z + region.Max.Z) / 2

		subregions := []Region{
			{Min: region.Min, Max: Point{midX, midY, midZ}},
			{Min: Point{midX, region.Min.Y, region.Min.Z}, Max: Point{region.Max.X, midY, midZ}},
			{Min: Point{region.Min.X, midY, region.Min.Z}, Max: Point{midX, region.Max.Y, midZ}},
			{Min: Point{midX, midY, region.Min.Z}, Max: Point{region.Max.X, region.Max.Y, midZ}},
			{Min: Point{region.Min.X, region.Min.Y, midZ}, Max: Point{midX, midY, region.Max.Z}},
			{Min: Point{midX, region.Min.Y, midZ}, Max: Point{region.Max.X, midY, region.Max.Z}},
			{Min: Point{region.Min.X, midY, midZ}, Max: Point{midX, region.Max.Y, region.Max.Z}},
			{Min: Point{midX, midY, midZ}, Max: region.Max},
		}

		// Count nanobots in range for each subregion and push to the queue
		for _, subregion := range subregions {
			count := 0
			for _, nanobot := range nanobots {
				if inRangeOfRegion(nanobot, subregion) {
					count++
				}
			}
			if count > 0 {
				heap.Push(&pq, &Region{Min: subregion.Min, Max: subregion.Max, Count: count})
			}
		}
	}

	return best
}

func manhattan(a, b Point) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y) + abs(a.Z-b.Z)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func inRangeOfRegion(nanobot Nanobot, region Region) bool {
	// Calculate the closest point in the region to the nanobot and check if it is in range
	closestPoint := Point{
		X: clamp(nanobot.X, region.Min.X, region.Max.X),
		Y: clamp(nanobot.Y, region.Min.Y, region.Max.Y),
		Z: clamp(nanobot.Z, region.Min.Z, region.Max.Z),
	}
	return manhattan(Point{nanobot.X, nanobot.Y, nanobot.Z}, closestPoint) <= nanobot.R
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
