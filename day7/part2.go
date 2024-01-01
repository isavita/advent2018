package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Step struct {
	id       rune
	duration int
}

func main() {
	// Read and parse the input file
	deps, allSteps := parseInput("day7/input.txt")
	// Simulate the work process
	timeTaken := simulateWork(deps, allSteps, 5, 60)
	fmt.Println(timeTaken)
}

func parseInput(filename string) (map[rune][]rune, map[rune]*Step) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	deps := make(map[rune][]rune)
	allSteps := make(map[rune]*Step)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var a, b rune
		fmt.Sscanf(scanner.Text(), "Step %c must be finished before step %c can begin.", &a, &b)
		deps[b] = append(deps[b], a)
		if _, exists := allSteps[a]; !exists {
			allSteps[a] = &Step{id: a, duration: int(a-'A') + 61}
		}
		if _, exists := allSteps[b]; !exists {
			allSteps[b] = &Step{id: b, duration: int(b-'A') + 61}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return deps, allSteps
}

func simulateWork(deps map[rune][]rune, allSteps map[rune]*Step, numWorkers int, baseDuration int) int {
	workers := make([]int, numWorkers)
	tasks := make([]rune, numWorkers)
	var time int

	for len(allSteps) > 0 {
		// Find available steps
		var available []rune
		for step, _ := range allSteps {
			if len(deps[step]) == 0 && !isBeingWorkedOn(step, tasks) {
				available = append(available, step)
			}
		}
		sort.Slice(available, func(i, j int) bool { return available[i] < available[j] })

		// Assign available steps to free workers
		for i := 0; i < numWorkers; i++ {
			if workers[i] == 0 && len(available) > 0 {
				tasks[i] = available[0]
				workers[i] = allSteps[available[0]].duration
				available = available[1:]
			}
		}

		// Advance time by the smallest non-zero duration among workers
		minDuration := findMinDuration(workers)
		for i := 0; i < numWorkers; i++ {
			if workers[i] != 0 {
				workers[i] -= minDuration
				if workers[i] == 0 {
					finishStep(deps, allSteps, tasks[i])
					tasks[i] = 0
				}
			}
		}
		time += minDuration
	}

	return time
}

func isBeingWorkedOn(step rune, tasks []rune) bool {
	for _, task := range tasks {
		if task == step {
			return true
		}
	}
	return false
}

func findMinDuration(durations []int) int {
	min := 1<<31 - 1
	for _, duration := range durations {
		if duration > 0 && duration < min {
			min = duration
		}
	}
	return min
}

func finishStep(deps map[rune][]rune, allSteps map[rune]*Step, step rune) {
	delete(allSteps, step)
	for s := range allSteps {
		deps[s] = remove(deps[s], step)
	}
}

func remove(slice []rune, val rune) []rune {
	for i, item := range slice {
		if item == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
