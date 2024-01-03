package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Cart struct {
	x, y      int
	direction rune
	turns     int
}

type Position struct {
	x, y int
}

func main() {
	file, _ := os.Open("day13/input.txt")
	defer file.Close()

	var tracks [][]rune
	var carts []*Cart
	scanner := bufio.NewScanner(file)

	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()
		trackLine := make([]rune, len(line))
		for x, r := range line {
			switch r {
			case '>', '<', '^', 'v':
				carts = append(carts, &Cart{x, y, r, 0})
				if r == '>' || r == '<' {
					trackLine[x] = '-'
				} else {
					trackLine[x] = '|'
				}
			default:
				trackLine[x] = r
			}
		}
		tracks = append(tracks, trackLine)
	}

	for len(carts) > 1 {
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y == carts[j].y {
				return carts[i].x < carts[j].x
			}
			return carts[i].y < carts[j].y
		})

		toRemove := make(map[int]struct{})
		for i, cart := range carts {
			if _, exists := toRemove[i]; exists {
				continue
			}

			moveCart(cart, tracks)
			if crashIndex := checkCrash(cart, carts); crashIndex != -1 {
				toRemove[i] = struct{}{}
				toRemove[crashIndex] = struct{}{}
			}
		}

		newCarts := []*Cart{}
		for i, cart := range carts {
			if _, exists := toRemove[i]; !exists {
				newCarts = append(newCarts, cart)
			}
		}
		carts = newCarts
	}

	fmt.Printf("%d,%d", carts[0].x, carts[0].y)
}

func moveCart(cart *Cart, tracks [][]rune) {
	switch cart.direction {
	case '>':
		cart.x++
	case '<':
		cart.x--
	case '^':
		cart.y--
	case 'v':
		cart.y++
	}

	switch tracks[cart.y][cart.x] {
	case '+':
		turnCart(cart)
	case '/', '\\':
		changeDirection(cart, tracks[cart.y][cart.x])
	}
}

func turnCart(cart *Cart) {
	if cart.turns%3 == 0 {
		if cart.direction == '>' {
			cart.direction = '^'

		} else if cart.direction == '<' {
			cart.direction = 'v'
		} else if cart.direction == '^' {
			cart.direction = '<'
		} else {
			cart.direction = '>'
		}
	} else if cart.turns%3 == 2 {
		if cart.direction == '>' {
			cart.direction = 'v'
		} else if cart.direction == '<' {
			cart.direction = '^'
		} else if cart.direction == '^' {
			cart.direction = '>'
		} else {
			cart.direction = '<'
		}
	}
	cart.turns++
}

func changeDirection(cart *Cart, track rune) {
	if track == '/' {
		if cart.direction == '>' {
			cart.direction = '^'
		} else if cart.direction == '<' {
			cart.direction = 'v'
		} else if cart.direction == '^' {
			cart.direction = '>'
		} else {
			cart.direction = '<'
		}
	} else if track == '\\' {
		if cart.direction == '>' {
			cart.direction = 'v'
		} else if cart.direction == '<' {
			cart.direction = '^'
		} else if cart.direction == '^' {
			cart.direction = '<'
		} else {
			cart.direction = '>'
		}
	}
}

func checkCrash(cart *Cart, carts []*Cart) int {
	for i, c := range carts {
		if c != cart && c.x == cart.x && c.y == cart.y {
			return i
		}
	}
	return -1
}
