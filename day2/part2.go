package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Open("day2/input.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	for i := 0; i < len(lines)-1; i++ {
		for j := i + 1; j < len(lines); j++ {
			diff := 0
			for k := 0; k < len(lines[i]); k++ {
				if lines[i][k] != lines[j][k] {
					diff++
					if diff > 1 {
						break
					}
				}
			}
			if diff == 1 {
				common := ""
				for k := 0; k < len(lines[i]); k++ {
					if lines[i][k] == lines[j][k] {
						common += string(lines[i][k])
					}
				}
				fmt.Println(common)
				return
			}
		}
	}
}
