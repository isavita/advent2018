package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Marble struct {
	value int
	prev  *Marble
	next  *Marble
}

func main() {
	players, lastMarble := readInput("day9/input.txt")
	fmt.Println(playMarbleGame(players, lastMarble))
}

func readInput(filename string) (int, int) {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	parts := strings.Fields(line)
	players, _ := strconv.Atoi(parts[0])
	lastMarble, _ := strconv.Atoi(parts[6])

	return players, lastMarble
}

func playMarbleGame(players, lastMarble int) int {
	scores := make([]int, players)
	current := &Marble{value: 0}
	current.next = current
	current.prev = current

	for marble := 1; marble <= lastMarble; marble++ {
		if marble%23 == 0 {
			player := marble % players
			for i := 0; i < 7; i++ {
				current = current.prev
			}
			scores[player] += marble + current.value
			current.prev.next = current.next
			current.next.prev = current.prev
			current = current.next
		} else {
			current = current.next
			newMarble := &Marble{value: marble, prev: current, next: current.next}
			current.next.prev = newMarble
			current.next = newMarble
			current = newMarble
		}
	}

	maxScore := 0
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}
