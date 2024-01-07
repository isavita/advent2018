package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type operation func(r []int, a, b int) int

var instructions = map[string]operation{
	"addr": func(r []int, a, b int) int { return r[a] + r[b] },
	"addi": func(r []int, a, b int) int { return r[a] + b },
	"mulr": func(r []int, a, b int) int { return r[a] * r[b] },
	"muli": func(r []int, a, b int) int { return r[a] * b },
	"banr": func(r []int, a, b int) int { return r[a] & r[b] },
	"bani": func(r []int, a, b int) int { return r[a] & b },
	"borr": func(r []int, a, b int) int { return r[a] | r[b] },
	"bori": func(r []int, a, b int) int { return r[a] | b },
	"setr": func(r []int, a, b int) int { return r[a] },
	"seti": func(r []int, a, b int) int { return a },
	"gtir": func(r []int, a, b int) int {
		if a > r[b] {
			return 1
		}
		return 0
	},
	"gtri": func(r []int, a, b int) int {
		if r[a] > b {
			return 1
		}
		return 0
	},
	"gtrr": func(r []int, a, b int) int {
		if r[a] > r[b] {
			return 1
		}
		return 0
	},
	"eqir": func(r []int, a, b int) int {
		if a == r[b] {
			return 1
		}
		return 0
	},
	"eqri": func(r []int, a, b int) int {
		if r[a] == b {
			return 1
		}
		return 0
	},
	"eqrr": func(r []int, a, b int) int {
		if r[a] == r[b] {
			return 1
		}
		return 0
	},
}

func loadProgram(lines []string) (int, []func([]int)) {
	var program []func([]int)
	var ipRegister int
	re := regexp.MustCompile(`\d+`)

	for _, line := range lines {
		if strings.HasPrefix(line, "#ip") {
			ipRegister, _ = strconv.Atoi(strings.Fields(line)[1])
			continue
		}

		parts := strings.Fields(line)
		op := instructions[parts[0]]
		nums := re.FindAllString(line, -1)
		a, _ := strconv.Atoi(nums[0])
		b, _ := strconv.Atoi(nums[1])
		c, _ := strconv.Atoi(nums[2])

		program = append(program, func(r []int) {
			r[c] = op(r, a, b)
		})
	}
	return ipRegister, program
}

func runProgram(ipRegister int, program []func([]int), registers []int, maxCycles int) []int {
	ip := 0
	cycles := 0

	for ip >= 0 && ip < len(program) {
		registers[ipRegister] = ip
		program[ip](registers)
		ip = registers[ipRegister] + 1
		cycles++
		if maxCycles > 0 && cycles >= maxCycles {
			break
		}
	}
	return registers
}

func max(slice []int) int {
	maxValue := slice[0]
	for _, v := range slice {
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

func main() {
	file, err := os.Open("day19/input.txt")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		lines = append(lines, text)
	}

	ipRegister, program := loadProgram(lines)

	registers := make([]int, 6)
	registers[0] = 1
	registers = runProgram(ipRegister, program, registers, 1000)
	n := max(registers)
	total := 0
	for i := 1; i <= n; i++ {
		if n%i == 0 {
			total += i
		}
	}
	fmt.Println(total)
}
