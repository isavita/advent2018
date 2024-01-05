package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("day19/input.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var ipBind int
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "#ip %d", &ipBind)

	var instructions [][]string
	for scanner.Scan() {
		instructions = append(instructions, strings.Fields(scanner.Text()))
	}

	registers := [6]int{}
	for ip := 0; ip < len(instructions); ip++ {
		registers[ipBind] = ip
		inst := instructions[ip]
		opcode, a, b, c := inst[0], toInt(inst[1]), toInt(inst[2]), toInt(inst[3])

		switch opcode {
		case "addr":
			registers[c] = registers[a] + registers[b]
		case "addi":
			registers[c] = registers[a] + b
		case "mulr":
			registers[c] = registers[a] * registers[b]
		case "muli":
			registers[c] = registers[a] * b
		case "banr":
			registers[c] = registers[a] & registers[b]
		case "bani":
			registers[c] = registers[a] & b
		case "borr":
			registers[c] = registers[a] | registers[b]
		case "bori":
			registers[c] = registers[a] | b
		case "setr":
			registers[c] = registers[a]
		case "seti":
			registers[c] = a
		case "gtir":
			registers[c] = boolToInt(a > registers[b])
		case "gtri":
			registers[c] = boolToInt(registers[a] > b)
		case "gtrr":
			registers[c] = boolToInt(registers[a] > registers[b])
		case "eqir":
			registers[c] = boolToInt(a == registers[b])
		case "eqri":
			registers[c] = boolToInt(registers[a] == b)
		case "eqrr":
			registers[c] = boolToInt(registers[a] == registers[b])
		}

		ip = registers[ipBind]
		if ip < 0 || ip >= len(instructions) {
			break
		}
	}

	fmt.Println(registers[0])
}

func toInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
