package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Record struct {
	timestamp time.Time
	action    string
	guardID   int
}

func main() {
	records := readAndParseInput("day4/input.txt")
	sort.Slice(records, func(i, j int) bool {
		return records[i].timestamp.Before(records[j].timestamp)
	})

	guardSleepMinutes := map[int][]int{}
	var currentGuardID int
	var sleepStart time.Time

	for _, record := range records {
		switch record.action {
		case "begins shift":
			currentGuardID = record.guardID
		case "falls asleep":
			sleepStart = record.timestamp
		case "wakes up":
			if _, ok := guardSleepMinutes[currentGuardID]; !ok {
				guardSleepMinutes[currentGuardID] = make([]int, 60)
			}
			for i := sleepStart.Minute(); i < record.timestamp.Minute(); i++ {
				guardSleepMinutes[currentGuardID][i]++
			}
		}
	}

	maxSleep := 0
	var sleepiestGuard int
	for guardID, minutes := range guardSleepMinutes {
		totalSleep := 0
		for _, count := range minutes {
			totalSleep += count
		}
		if totalSleep > maxSleep {
			maxSleep = totalSleep
			sleepiestGuard = guardID
		}
	}

	maxMinute := 0
	maxMinuteCount := 0
	for i, count := range guardSleepMinutes[sleepiestGuard] {
		if count > maxMinuteCount {
			maxMinuteCount = count
			maxMinute = i
		}
	}

	fmt.Println(sleepiestGuard * maxMinute)
}

func readAndParseInput(filename string) []Record {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	layout := "2006-01-02 15:04"
	var records []Record

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "] ")
		timePart := parts[0][1:]
		actionPart := parts[1]

		ts, err := time.Parse(layout, timePart)
		if err != nil {
			panic(err)
		}

		guardID := -1
		if strings.Contains(actionPart, "Guard") {
			fmt.Sscanf(actionPart, "Guard #%d begins shift", &guardID)
			actionPart = "begins shift"
		} else if strings.Contains(actionPart, "falls asleep") {
			actionPart = "falls asleep"
		} else if strings.Contains(actionPart, "wakes up") {
			actionPart = "wakes up"
		}

		records = append(records, Record{timestamp: ts, action: actionPart, guardID: guardID})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return records
}
