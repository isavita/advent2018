package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Read input from file
	numbers := readInput("day8/input.txt")
	// Parse the tree and calculate the value of the root node
	value, _ := parseTree(numbers, 0)
	fmt.Println(value)
}

func readInput(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	parts := strings.Split(line, " ")
	numbers := make([]int, len(parts))
	for i, part := range parts {
		numbers[i], err = strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
	}

	return numbers
}

func parseTree(data []int, index int) (int, int) {
	childCount, metaCount := data[index], data[index+1]
	index += 2

	childValues := make([]int, childCount)
	for i := 0; i < childCount; i++ {
		var childValue int
		childValue, index = parseTree(data, index)
		childValues[i] = childValue
	}

	value := 0
	if childCount == 0 {
		// If no children, sum the metadata entries
		for i := 0; i < metaCount; i++ {
			value += data[index+i]
		}
	} else {
		// If there are children, sum the values of the referenced children
		for i := 0; i < metaCount; i++ {
			metadata := data[index+i]
			if metadata <= childCount && metadata > 0 {
				value += childValues[metadata-1]
			}
		}
	}
	index += metaCount

	return value, index
}
