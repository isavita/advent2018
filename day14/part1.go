package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, _ := os.Open("day14/input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	input, _ := strconv.Atoi(scanner.Text())

	scoreboard := []int{3, 7}
	elf1, elf2 := 0, 1

	for len(scoreboard) < input+10 {
		newScore := scoreboard[elf1] + scoreboard[elf2]
		if newScore >= 10 {
			scoreboard = append(scoreboard, newScore/10)
		}
		scoreboard = append(scoreboard, newScore%10)

		elf1 = (elf1 + scoreboard[elf1] + 1) % len(scoreboard)
		elf2 = (elf2 + scoreboard[elf2] + 1) % len(scoreboard)
	}

	for i := input; i < input+10; i++ {
		fmt.Print(scoreboard[i])
	}
	fmt.Println()
}
