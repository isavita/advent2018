package main

import (
	"container/heap"
	"fmt"
	"os"
)

type coord struct {
	x, y int
}

func (c coord) neighbors() []coord {
	n := []coord{
		{c.x + 1, c.y},
		{c.x, c.y + 1},
	}

	if c.x > 0 {
		n = append(n, coord{c.x - 1, c.y})
	}
	if c.y > 0 {
		n = append(n, coord{c.x, c.y - 1})
	}

	return n
}

const (
	geologicY  = 16807
	geologicX  = 48271
	caveModulo = 20183
)

const (
	TypeRocky  = 0
	TypeWet    = 1
	TypeNarrow = 2
)

const (
	ToolNone = 1 << iota
	ToolTorch
	ToolGear
)

type Map struct {
	target               coord
	depth                int
	geologicIndicesCache map[int]map[int]int
	erosionLevelsCache   map[int]map[int]int
}

func NewMap(input string) *Map {
	m := Map{}
	m.geologicIndicesCache = make(map[int]map[int]int)
	m.erosionLevelsCache = make(map[int]map[int]int)

	_, err := fmt.Sscanf(input, "depth: %d\ntarget: %d,%d", &m.depth, &m.target.x, &m.target.y)
	if err != nil {
		panic(err)
	}

	return &m
}

func (m *Map) GeologicIndex(x, y int) int {
	if m.geologicIndicesCache[y] != nil {
		if i, ok := m.geologicIndicesCache[y][x]; ok {
			return i
		}
	} else {
		m.geologicIndicesCache[y] = make(map[int]int)
	}

	switch {
	case x == 0 && y == 0, x == m.target.x && y == m.target.y:
		m.geologicIndicesCache[y][x] = 0
	case y == 0:
		m.geologicIndicesCache[y][x] = x * geologicY
	case x == 0:
		m.geologicIndicesCache[y][x] = y * geologicX
	default:
		m.geologicIndicesCache[y][x] = m.ErosionLevel(x-1, y) * m.ErosionLevel(x, y-1)
	}

	return m.geologicIndicesCache[y][x]
}

func (m *Map) ErosionLevel(x, y int) int {
	if m.erosionLevelsCache[y] != nil {
		if level, ok := m.erosionLevelsCache[y][x]; ok {
			return level
		}
	} else {
		m.erosionLevelsCache[y] = make(map[int]int)
	}

	m.erosionLevelsCache[y][x] = (m.GeologicIndex(x, y) + m.depth) % caveModulo

	return m.erosionLevelsCache[y][x]
}

func (m Map) Type(x, y int) int {
	return m.ErosionLevel(x, y) % 3
}

func (m Map) Neighbors(pos coord, equip int) []Item {
	var n []Item

	for _, c := range pos.neighbors() {
		t := m.Type(c.x, c.y)

		if equip&allowed(t) != 0 {
			n = append(n, Item{pos: c, equip: equip, time: 1})
			n = append(n, Item{pos: c, equip: equip ^ allowed(t), time: 8})
		}
	}

	return n
}

func allowed(regionType int) int {
	switch regionType {
	case TypeRocky:
		return ToolGear | ToolTorch
	case TypeWet:
		return ToolGear | ToolNone
	case TypeNarrow:
		return ToolTorch | ToolNone
	default:
		panic(fmt.Errorf("unknown region type: %d", regionType))
	}
}

type Item struct {
	pos   coord
	equip int
	time  int
	index int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].time < pq[j].time
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.index = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (i Item) Compare(o Item) int {
	switch {
	case i.time > o.time:
		return 1
	case i.time < o.time:
		return -1
	}
	return 0
}

const bailFactor = 8

func Rescue(input string) int {
	m := NewMap(input)

	queue := PriorityQueue{
		&Item{pos: coord{0, 0}, time: 0, equip: ToolTorch},
	}
	heap.Init(&queue)

	type step struct {
		coord coord
		equip int
	}

	distances := map[step]int{
		step{coord: coord{0, 0}, equip: ToolTorch}: 0,
	}

	for len(queue) > 0 {
		item := (heap.Pop(&queue)).(*Item)

		if item.pos.x == m.target.x && item.pos.y == m.target.y && item.equip == ToolTorch {
			return item.time
		}

		if item.pos.x > bailFactor*m.target.x || item.pos.y > bailFactor*m.target.y {
			// don't wander too far off
			continue
		}

		if t, ok := distances[step{coord: item.pos, equip: item.equip}]; ok && t < item.time {
			continue
		}

		for _, n := range m.Neighbors(item.pos, item.equip) {
			d := step{coord: n.pos, equip: n.equip}

			if t, ok := distances[step{coord: n.pos, equip: n.equip}]; !ok || item.time+n.time < t {
				distances[d] = item.time + n.time
				heap.Push(&queue, &Item{pos: n.pos, time: item.time + n.time, equip: n.equip})
			}
		}
	}

	return 0
}

func main() {
	input, err := os.ReadFile("day22/input.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(Rescue(string(input)))
}
