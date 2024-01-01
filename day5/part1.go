package main

import (
	"bufio"
	"fmt"
	"os"
)

func react(polymer string) string {
	for i := 0; i < len(polymer)-1; i++ {
		if polymer[i] != polymer[i+1] &&
			(polymer[i]+32 == polymer[i+1] ||
				polymer[i]-32 == polymer[i+1]) {
			return react(polymer[:i] + polymer[i+2:])
		}
	}
	return polymer
}

func main() {
	file, _ := os.Open("day5/input.txt")
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	polymer := scanner.Text()

	result := react(polymer)
	fmt.Println(len(result))
}
