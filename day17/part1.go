package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("day17/input.txt")
	if err != nil {
		panic(err)
	}
	inputStr := strings.TrimSpace(string(input))
	lines := regSplit(inputStr, "\n")

	ground := make([][]rune, 1)
	ground[0] = make([]rune, 1)
	ground[0][0] = '+'

	maxX, minX, maxY, minY := 0, 0, 0, 20
	xOffset, yOffset := 500, 0

	for _, line := range lines {
		split := regSplit(line, "[=, .]+")
		if split[0] == "x" {
			x := strToInt(split[1]) - xOffset
			y1 := strToInt(split[3]) - yOffset
			y2 := strToInt(split[4]) - yOffset

			for x >= maxX {
				maxX++
				for j := range ground {
					ground[j] = append(ground[j], '.')
				}
			}
			for x <= minX {
				minX--
				for j := range ground {
					ground[j] = append([]rune{'.'}, ground[j]...)
				}
			}
			for y2 > maxY {
				maxY++
				ground = append(ground, make([]rune, len(ground[0])))
				for j := range ground[len(ground)-1] {
					ground[len(ground)-1][j] = '.'
				}
			}
			if y1 < minY {
				minY = y1
			}
			for i := y1; i <= y2; i++ {
				ground[i][x-minX] = '#'
			}

		} else {
			y := strToInt(split[1]) - yOffset
			x1 := strToInt(split[3]) - xOffset
			x2 := strToInt(split[4]) - xOffset

			for y > maxY {
				maxY++
				ground = append(ground, make([]rune, len(ground[0])))
				for j := range ground[len(ground)-1] {
					ground[len(ground)-1][j] = '.'
				}
			}
			for x2 >= maxX {
				maxX++
				for j := range ground {
					ground[j] = append(ground[j], '.')
				}
			}
			for x1 <= minX {
				minX--
				for j := range ground {
					ground[j] = append([]rune{'.'}, ground[j]...)
				}
			}
			for i := x1; i <= x2; i++ {
				ground[y][i-minX] = '#'
			}
			if y < minY {
				minY = y
			}
		}
	}

	waterCount := 0
	flowCount := 0
	roundLimit := 200000

	for ground[1][-minX] != '|' && waterCount < roundLimit {
		canMove := true
		x := -minX
		y := 1
		tryLeft := 0
		for canMove {
			if y+1 > maxY || ground[y+1][x] == '|' {
				ground[y][x] = '|'
				canMove = false
				if y >= minY {
					flowCount++
				}
			} else if ground[y+1][x] == '.' {
				y++
				tryLeft = 0
			} else if ground[y+1][x] == '#' || ground[y+1][x] == '~' {
				if (tryLeft == 1 && ground[y][x-1] == '|') ||
					(tryLeft == 2 && ground[y][x+1] == '|') ||
					(ground[y][x+1] == '|' && ground[y][x-1] != '.') ||
					(ground[y][x+1] != '.' && ground[y][x-1] == '|') {
					ground[y][x] = '|'
					flowCount++
					canMove = false
					for i := x + 1; ground[y][i] == '~'; i++ {
						ground[y][i] = '|'
						waterCount--
						flowCount++
					}
					for i := x - 1; ground[y][i] == '~'; i-- {
						ground[y][i] = '|'
						waterCount--
						flowCount++
					}
				} else if (tryLeft == 0 && ground[y][x-1] == '.') ||
					(tryLeft == 1 && ground[y][x-1] == '.') {
					x--
					tryLeft = 1
				} else if (tryLeft == 0 && ground[y][x+1] == '.') ||
					(tryLeft == 2 && ground[y][x+1] == '.') {
					x++
					tryLeft = 2
				} else {
					canMove = false
					ground[y][x] = '~'
					waterCount++
				}
			}

		}

	}
	fmt.Println(flowCount + waterCount)
}

func strToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	i := reg.FindAllStringIndex(text, -1)
	laststart := 0
	n := make([]string, len(i)+1)
	for i, element := range i {
		n[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	n[len(i)] = text[laststart:]
	return n
}
