package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Claim represents a single fabric claim by an elf
type Claim struct {
	id     int
	left   int
	top    int
	width  int
	height int
}

// ParseClaim parses a string into a Claim struct
func ParseClaim(s string) (Claim, error) {
	var c Claim
	_, err := fmt.Sscanf(s, "#%d @ %d,%d: %dx%d",
		&c.id, &c.left, &c.top, &c.width, &c.height)
	return c, err
}

// ReadClaims reads claims from a file
func ReadClaims(filename string) ([]Claim, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var claims []Claim
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		claim, err := ParseClaim(line)
		if err != nil {
			return claims, err
		}
		claims = append(claims, claim)
	}
	return claims, scanner.Err()
}

// CountOverlappingInches counts the number of square inches that overlap
func CountOverlappingInches(claims []Claim) int {
	fabric := make(map[string]int)
	for _, claim := range claims {
		for i := claim.left; i < claim.left+claim.width; i++ {
			for j := claim.top; j < claim.top+claim.height; j++ {
				coord := strconv.Itoa(i) + "," + strconv.Itoa(j)
				fabric[coord]++
			}
		}
	}

	overlapping := 0
	for _, count := range fabric {
		if count > 1 {
			overlapping++
		}
	}
	return overlapping
}

func main() {
	claims, err := ReadClaims("day3/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	overlapping := CountOverlappingInches(claims)
	fmt.Println(overlapping)
}
