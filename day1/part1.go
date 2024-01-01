package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	freqChanges := readInput()
	freq := 0
	for _, change := range freqChanges {
		freq += parseChange(change)
	}
	fmt.Println(freq)
}

func readInput() []string {
	var lines []string
	file, err := openFile("day1/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return lines
}

func openFile(path string) (*os.File, error) {
	return os.Open(path)
}

func parseChange(change string) int {
	sign, num, err := parseSignNum(change)
	if err != nil {
		panic(err)
	}
	return sign * num
}

func parseSignNum(change string) (int, int, error) {
	sign := 1
	if change[0] == '-' {
		sign = -1
		change = change[1:]
	}
	num, err := strconv.Atoi(change)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid frequency change: %s", change)
	}
	return sign, num, nil
}
