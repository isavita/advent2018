package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("day2/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	twoCount, threeCount := 0, 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		id := scanner.Text()
		twos, threes := countTwosAndThrees(id)
		if twos {
			twoCount++
		}
		if threes {
			threeCount++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	checksum := twoCount * threeCount
	fmt.Println(checksum)
}

func countTwosAndThrees(id string) (bool, bool) {
	charCount := make(map[rune]int)
	for _, char := range id {
		charCount[char]++
	}

	var hasTwos, hasThrees bool
	for _, count := range charCount {
		if count == 2 {
			hasTwos = true
		} else if count == 3 {
			hasThrees = true
		}
	}
	return hasTwos, hasThrees
}
