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

	for generation := 0; generation < 20; generation++ {
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
	}

	sum := 0
	for k := range state {
		sum += k
	}

	fmt.Println(sum)
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
