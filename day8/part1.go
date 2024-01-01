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
	// Parse the tree and calculate the metadata sum
	sum, _ := parseTree(numbers, 0)
	fmt.Println(sum)
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

	sum := 0
	for i := 0; i < childCount; i++ {
		childSum, newIndex := parseTree(data, index)
		sum += childSum
		index = newIndex
	}

	for i := 0; i < metaCount; i++ {
		sum += data[index+i]
	}
	index += metaCount

	return sum, index
}
