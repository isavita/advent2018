package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Claim struct {
	ID     int
	X      int
	Y      int
	Width  int
	Height int
}

func main() {
	claims, err := readClaims("day3/input.txt")
	if err != nil {
		log.Fatalf("read claims: %v", err)
	}

	fabric := make([][]int, 1000)
	for i := range fabric {
		fabric[i] = make([]int, 1000)
	}

	for _, claim := range claims {
		for y := claim.Y; y < claim.Y+claim.Height; y++ {
			for x := claim.X; x < claim.X+claim.Width; x++ {
				fabric[y][x]++
			}
		}
	}

	for _, claim := range claims {
		overlap := false
		for y := claim.Y; y < claim.Y+claim.Height; y++ {
			for x := claim.X; x < claim.X+claim.Width; x++ {
				if fabric[y][x] > 1 {
					overlap = true
					break
				}
			}
			if overlap {
				break
			}
		}
		if !overlap {
			fmt.Println(claim.ID)
			return
		}
	}
}

func readClaims(filename string) ([]Claim, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var claims []Claim
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		id, _ := strconv.Atoi(parts[0][1:])
		coords := strings.Split(parts[2][:len(parts[2])-1], ",")
		x, _ := strconv.Atoi(coords[0])
		y, _ := strconv.Atoi(coords[1])
		dims := strings.Split(parts[3], "x")
		width, _ := strconv.Atoi(dims[0])
		height, _ := strconv.Atoi(dims[1])
		claims = append(claims, Claim{id, x, y, width, height})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return claims, nil
}
