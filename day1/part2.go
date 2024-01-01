package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile("day1/input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	frequencyChanges := strings.Split(strings.TrimSpace(string(data)), "\n")
	frequencies := make(map[int]bool)
	currentFrequency := 0
	frequencies[currentFrequency] = true

	for {
		for _, change := range frequencyChanges {
			frequencyDelta, err := strconv.Atoi(change)
			if err != nil {
				fmt.Println("Error converting string to int:", err)
				os.Exit(1)
			}
			currentFrequency += frequencyDelta
			if frequencies[currentFrequency] {
				fmt.Println(currentFrequency)
				return
			}
			frequencies[currentFrequency] = true
		}
	}
}
