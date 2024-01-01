package main

import (
	"fmt"
	"os"
	"strings"
)

func react(polymer string) string {
	reactionOccurred := true
	for reactionOccurred {
		reactionOccurred = false
		for i := 0; i < len(polymer)-1; i++ {
			if polymer[i] != polymer[i+1] && strings.ToUpper(string(polymer[i])) == strings.ToUpper(string(polymer[i+1])) {
				polymer = polymer[:i] + polymer[i+2:]
				reactionOccurred = true
			}
		}
	}
	return polymer
}

func main() {
	content, err := os.ReadFile("day5/input.txt")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	polymer := string(content)

	minLength := len(polymer)
	for unit := 'a'; unit <= 'z'; unit++ {
		tempPolymer := strings.ReplaceAll(polymer, string(unit), "")
		tempPolymer = strings.ReplaceAll(tempPolymer, strings.ToUpper(string(unit)), "")
		reactedPolymer := react(tempPolymer)
		if len(reactedPolymer) < minLength {
			minLength = len(reactedPolymer)
		}
	}

	fmt.Println(minLength - 1)
}
