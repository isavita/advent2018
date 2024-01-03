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
	input := scanner.Text()

	scoreboard := []int{3, 7}
	elf1, elf2 := 0, 1
	inputLen := len(input)
	inputSequence := make([]int, inputLen)

	for i := 0; i < inputLen; i++ {
		inputSequence[i], _ = strconv.Atoi(string(input[i]))
	}

	for {
		newScore := scoreboard[elf1] + scoreboard[elf2]
		if newScore >= 10 {
			scoreboard = append(scoreboard, newScore/10)
			if checkSequence(scoreboard, inputSequence) {
				break
			}
		}
		scoreboard = append(scoreboard, newScore%10)
		if checkSequence(scoreboard, inputSequence) {
			break
		}

		elf1 = (elf1 + scoreboard[elf1] + 1) % len(scoreboard)
		elf2 = (elf2 + scoreboard[elf2] + 1) % len(scoreboard)
	}

	fmt.Println(len(scoreboard) - inputLen)
}

func checkSequence(scoreboard, sequence []int) bool {
	if len(scoreboard) < len(sequence) {
		return false
	}
	start := len(scoreboard) - len(sequence)
	for i, v := range sequence {
		if scoreboard[start+i] != v {
			return false
		}
	}
	return true
}
