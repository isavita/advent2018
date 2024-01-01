package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// Read and parse the input file
	deps, allSteps := parseInput("day7/input.txt")
	// Compute the order of steps
	order := topologicalSort(deps, allSteps)
	fmt.Println(order)
}

func parseInput(filename string) (map[rune][]rune, map[rune]bool) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	deps := make(map[rune][]rune)
	allSteps := make(map[rune]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var a, b rune
		fmt.Sscanf(scanner.Text(), "Step %c must be finished before step %c can begin.", &a, &b)
		deps[b] = append(deps[b], a)
		allSteps[a] = true
		allSteps[b] = true
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return deps, allSteps
}

func topologicalSort(deps map[rune][]rune, allSteps map[rune]bool) string {
	var order []rune
	available := make([]rune, 0)

	// Find initial available steps (with no dependencies)
	for step := range allSteps {
		if len(deps[step]) == 0 {
			available = append(available, step)
		}
	}
	sort.Slice(available, func(i, j int) bool { return available[i] < available[j] })

	for len(available) > 0 {
		// Pick the step that comes first alphabetically
		next := available[0]
		available = available[1:]
		order = append(order, next)

		// Check and update the availability of the next steps
		for step := range allSteps {
			if contains(deps[step], next) {
				deps[step] = remove(deps[step], next)
				if len(deps[step]) == 0 {
					available = append(available, step)
				}
			}
		}
		sort.Slice(available, func(i, j int) bool { return available[i] < available[j] })
	}
	return string(order)
}

func contains(slice []rune, val rune) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func remove(slice []rune, val rune) []rune {
	for i, item := range slice {
		if item == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
