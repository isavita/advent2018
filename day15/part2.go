package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Cave struct {
	Units SortableUnits
	Map   Map
}

const (
	KindSpace = 1 << iota
	KindElf
	KindGoblin
	KindWall
	KindHighlight
)

var KindRunes = map[int]rune{
	KindSpace:     '.',
	KindElf:       'E',
	KindGoblin:    'G',
	KindWall:      '#',
	KindHighlight: '@',
}

var RuneKinds = map[rune]int{
	'.': KindSpace,
	'E': KindElf,
	'G': KindGoblin,
	'#': KindWall,
}

func IsUnit(bit int) bool {
	return (KindElf|KindGoblin)&bit != 0
}

func NewCave(input []string, elfPower int) *Cave {
	c := &Cave{}
	c.ParseMap(input, elfPower)
	return c
}

func (c *Cave) ParseMap(input []string, elfPower int) {
	m := make(Map)

	for y, row := range input {
		for x, col := range row {
			kind, ok := RuneKinds[col]
			if !ok {
				kind = KindWall
			}

			tile := &Tile{Kind: kind}
			if IsUnit(kind) {
				c.Units = append(c.Units, NewUnit(tile, kind, elfPower))
			}
			m.SetTile(tile, x, y)
		}
	}
	c.Map = m
}

func (c Cave) PrintMap(highlight *Tile) {
	for y := 0; y < len(c.Map); y++ {
		var units []string
		for x := 0; x < len(c.Map[y]); x++ {
			t := c.Map.Tile(x, y)
			if t == highlight {
				fmt.Print(string(KindRunes[KindHighlight]))
			} else {
				fmt.Print(string(KindRunes[t.Kind]))
			}

			if t.Unit != nil {
				units = append(units, fmt.Sprintf("%c(%d)", KindRunes[t.Unit.Kind], t.Unit.Hitpoints))
			}
		}
		if len(units) > 0 {
			fmt.Print("  ", strings.Join(units, ", "))
		}
		fmt.Println()
	}
}

func (c Cave) PrintDistance(t *Tile) {
	distances, _ := c.Map.FindWalkableTiles(t)
	for y := 0; y < len(c.Map); y++ {
		for x := 0; x < len(c.Map[y]); x++ {
			curT := c.Map.Tile(x, y)
			if d, ok := distances[curT]; ok && curT != t {
				fmt.Print(d)
			} else {
				fmt.Print(string(KindRunes[curT.Kind]))
			}
		}
		fmt.Println()
	}
}

// Status returns the sum of remaining hitpoints for all units and true if the fight is still ongoing.
func (c Cave) Status() (int, bool) {
	var elves, goblins bool
	var hp int

	for _, u := range c.Units {
		if u.Hitpoints <= 0 {
			continue
		}
		if u.Kind == KindElf {
			elves = true
		} else {
			goblins = true
		}
		hp += u.Hitpoints
	}

	return hp, elves && goblins
}

func (c *Cave) RemoveTheDead() {
	var newUnits SortableUnits
	for _, unit := range c.Units {
		if unit.Hitpoints > 0 {
			newUnits = append(newUnits, unit)
		}
	}
	c.Units = newUnits
}

func (c *Cave) RemoveUnit(u *Unit) {
	u.Tile.Kind = KindSpace
	u.Tile.Unit = nil
	u.Tile = nil
}

// Tick returns false if combat ended during the round, and whether or not an elf has died this round.
func (c *Cave) Tick(stopOnElfDeath bool) (bool, bool) {
	c.RemoveTheDead()
	sort.Sort(c.Units)

	for _, unit := range c.Units {
		if unit.Hitpoints <= 0 {
			continue
		}
		if !unit.Targets(c) {
			return false, false
		}
		unit.Move(c)
		if unit.Attack(c) && stopOnElfDeath {
			return false, true
		}
	}
	return true, false
}

type Coordinate struct {
	X, Y int
}

var offsets = []Coordinate{
	{0, -1},
	{-1, 0},
	{1, 0},
	{0, 1},
}

type Map map[int]map[int]*Tile

func (m Map) SetTile(t *Tile, x, y int) {
	if m[y] == nil {
		m[y] = make(map[int]*Tile)
	}
	m[y][x] = t
	t.X = x
	t.Y = y
	t.Map = m
}

func (m Map) Tile(x, y int) *Tile {
	if m[y] == nil {
		return nil
	}
	return m[y][x]
}

func (m Map) FindWalkableTiles(t *Tile) (map[*Tile]int, map[*Tile]*Tile) {
	frontier := []*Tile{t}
	distance := map[*Tile]int{t: 0}
	cameFrom := map[*Tile]*Tile{t: nil}

	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]

		for _, next := range current.WalkableNeighbors() {
			if _, ok := distance[next]; !ok {
				frontier = append(frontier, next)
				distance[next] = distance[current] + 1
				cameFrom[next] = current
			}
		}
	}

	return distance, cameFrom
}

type Tile struct {
	Kind int
	X, Y int
	Map  Map
	Unit *Unit
}

func (t Tile) WalkableNeighbors() []*Tile {
	var neighbors []*Tile

	for _, offset := range offsets {
		if n := t.Map.Tile(t.X+offset.X, t.Y+offset.Y); n != nil && n.Kind == KindSpace {
			neighbors = append(neighbors, n)
		}
	}

	return neighbors
}

type SortableTiles []*Tile

func (s SortableTiles) Len() int {
	return len(s)
}

func (s SortableTiles) Less(i, j int) bool {
	if s[i].Y == s[j].Y {
		return s[i].X < s[j].X
	}
	return s[i].Y < s[j].Y
}

func (s SortableTiles) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type Unit struct {
	Kind      int
	Hitpoints int
	Power     int
	Tile      *Tile
}

const (
	defaultHitpoints = 200
	defaultPower     = 3
)

type SortableUnits []*Unit

func (s SortableUnits) Len() int {
	return len(s)
}

func (s SortableUnits) Less(i, j int) bool {
	if s[i].Tile.Y == s[j].Tile.Y {
		return s[i].Tile.X < s[j].Tile.X
	}
	return s[i].Tile.Y < s[j].Tile.Y
}

func (s SortableUnits) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func NewUnit(tile *Tile, kind, elfPower int) *Unit {
	unit := &Unit{
		Kind:      kind,
		Hitpoints: defaultHitpoints,
		Power:     defaultPower,
		Tile:      tile,
	}
	tile.Unit = unit
	if unit.Kind == KindElf {
		unit.Power = elfPower
	}
	return unit
}

func (u Unit) Targets(c *Cave) bool {
	for _, unit := range c.Units {
		if unit.Kind != u.Kind && unit.Hitpoints > 0 {
			return true
		}
	}
	return false
}

// NextTile returns the next tile a unit should move to and the target tile it is moving towards, or nil if no reachable
// target has been found.
func (u *Unit) NextTile(c *Cave) (*Tile, *Tile) {
	var targets SortableTiles

	closestTargetDistance := math.MaxInt32
	distances, path := c.Map.FindWalkableTiles(u.Tile)
	enemies := u.Enemies(c)

	for _, enemy := range enemies {
		for _, target := range enemy.Tile.WalkableNeighbors() {
			if distance, ok := distances[target]; ok && distance <= closestTargetDistance {
				if distance < closestTargetDistance {
					closestTargetDistance = distance
					targets = SortableTiles{}
				}
				targets = append(targets, target)
			}
		}
	}
	sort.Sort(targets)
	if len(targets) > 0 {
		target := targets[0]
		current := target
		for {
			if path[current] == u.Tile {
				return current, target
			}
			current = path[current]
		}
	}
	return nil, nil
}

// Enemies returns a list of enemy units sorted by map position in reading order
func (u *Unit) Enemies(c *Cave) SortableUnits {
	var enemies SortableUnits
	for _, unit := range c.Units {
		if unit.Kind != u.Kind && unit.Hitpoints > 0 {
			enemies = append(enemies, unit)
		}
	}
	sort.Sort(enemies)
	return enemies
}

func (u *Unit) EnemyNeighbor(c *Cave) *Unit {
	var target *Unit
	for _, offset := range offsets {
		if t := c.Map.Tile(u.Tile.X+offset.X, u.Tile.Y+offset.Y); t != nil && t.Unit != nil && t.Unit.Kind != u.Kind && t.Unit.Hitpoints > 0 {
			if target == nil || t.Unit.Hitpoints < target.Hitpoints {
				target = t.Unit
			}
		}
	}
	return target
}

func (u *Unit) Move(c *Cave) {
	if u.EnemyNeighbor(c) != nil {
		return
	}
	next, _ := u.NextTile(c)
	if next != nil {
		next.Unit = u
		next.Kind = u.Kind
		u.Tile.Kind = KindSpace
		u.Tile.Unit = nil
		u.Tile = next
	}
}

func (u *Unit) Attack(c *Cave) bool {
	enemy := u.EnemyNeighbor(c)
	if enemy != nil {
		killed := enemy.Damage(c, u.Power)
		return killed && enemy.Kind == KindElf
	}
	return false
}

func (u *Unit) Damage(c *Cave, damage int) bool {
	u.Hitpoints -= damage
	if u.Hitpoints <= 0 {
		c.RemoveUnit(u)
		return true
	}
	return false
}

func Combat(input []string) int {
	cave := NewCave(input, defaultPower)
	for i := 1; true; i++ {
		hp, combat := cave.Status()

		if !combat {
			return (i - 1) * hp
		}

		if cleanRound, _ := cave.Tick(false); !cleanRound {
			i--
		}
	}
	return -1
}

func CheatingElves(input []string) int {
	elfDied := true
	for power := 4; elfDied; power++ {
		cave := NewCave(input, power)
		for i := 1; true; i++ {
			hp, combat := cave.Status()

			if !combat {
				return (i - 1) * hp
			}

			var cleanRound bool
			cleanRound, elfDied = cave.Tick(true)
			if elfDied {
				break
			}
			if !cleanRound {
				i--
			}
		}
	}
	return -1
}

func main() {
	input, err := os.ReadFile("day15/input.txt")
	if err != nil {
		panic(err)
	}

	lines := strings.Split(strings.TrimSpace(string(input)), "\n")
	fmt.Println(CheatingElves(lines))
}
