package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("day12/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var initialState string
	rules := make(map[string]byte)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "initial state") {
			initialState = strings.Split(line, ": ")[1]
		} else if strings.Contains(line, "=>") {
			parts := strings.Split(line, " => ")
			rules[parts[0]] = parts[1][0]
		}
	}

	state := make(map[int]byte)
	for i, c := range initialState {
		if c == '#' {
			state[i] = '#'
		}
	}

	previousPattern := ""
	previousSum := 0
	offset := 0
	for generation := 0; generation < 50000000000; generation++ {
		newState := make(map[int]byte)
		minPot, maxPot := minMaxKeys(state)
		for i := minPot - 2; i <= maxPot+2; i++ {
			pattern := ""
			for j := i - 2; j <= i+2; j++ {
				if state[j] == '#' {
					pattern += "#"
				} else {
					pattern += "."
				}
			}
			if rules[pattern] == '#' {
				newState[i] = '#'
			}
		}
		state = newState

		currentPattern, currentSum := statePattern(state)
		if currentPattern == previousPattern {
			// Once a repeating pattern is found, calculate the offset per generation
			offset = currentSum - previousSum
			remainingGenerations := 50000000000 - generation - 1
			finalSum := currentSum + offset*remainingGenerations
			fmt.Println(finalSum)
			return
		}
		previousPattern = currentPattern
		previousSum = currentSum
	}
}

func minMaxKeys(m map[int]byte) (minKey int, maxKey int) {
	first := true
	for k := range m {
		if first {
			minKey, maxKey = k, k
			first = false
		} else {
			if k < minKey {
				minKey = k
			}
			if k > maxKey {
				maxKey = k
			}
		}
	}
	return
}

func statePattern(m map[int]byte) (string, int) {
	minPot, maxPot := minMaxKeys(m)
	var pattern strings.Builder
	sum := 0
	for i := minPot; i <= maxPot; i++ {
		if m[i] == '#' {
			pattern.WriteByte('#')
			sum += i
		} else {
			pattern.WriteByte('.')
		}
	}
	return pattern.String(), sum
}
