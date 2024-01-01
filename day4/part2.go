package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type record struct {
	time  time.Time
	event string
}

type guard struct {
	id       int
	minutes  [60]int
	totalMin int
}

func main() {
	inputFile, _ := os.Open("day4/input.txt")
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	records := []record{}
	guards := make(map[int]*guard)

	for scanner.Scan() {
		line := scanner.Text()
		t, _ := time.Parse("2006-01-02 15:04", line[1:17])
		records = append(records, record{time: t, event: line[19:]})
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].time.Before(records[j].time)
	})

	var currentGuard *guard
	var sleepStart int

	for _, record := range records {
		switch {
		case strings.Contains(record.event, "begins shift"):
			id, _ := strconv.Atoi(strings.Split(record.event, " ")[1][1:])
			if _, ok := guards[id]; !ok {
				guards[id] = &guard{id: id}
			}
			currentGuard = guards[id]
		case strings.Contains(record.event, "falls asleep"):
			sleepStart = record.time.Minute()
		case strings.Contains(record.event, "wakes up"):
			for i := sleepStart; i < record.time.Minute(); i++ {
				currentGuard.minutes[i]++
				currentGuard.totalMin++
			}
		}
	}

	var mostFreqGuard *guard
	var mostFreqMin int

	for _, g := range guards {
		for i, m := range g.minutes {
			if mostFreqGuard == nil || m > mostFreqGuard.minutes[mostFreqMin] {
				mostFreqGuard = g
				mostFreqMin = i
			}
		}
	}

	fmt.Println(mostFreqGuard.id * mostFreqMin)
}
