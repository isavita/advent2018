package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.ReadFile("day21/input.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(solve(string(file)))
}

func solve(input string) int {
	opcodeComputer := parseInput(input)

	for !opcodeComputer.tick() {

		if opcodeComputer.registers[opcodeComputer.instructionPointer] == 28 {
			break
		}
	}

	return opcodeComputer.registers[5]
}

type opcodeComputer struct {
	instructions       []instruction
	registers          [6]int
	instructionPointer int
}
type instruction struct {
	name      string
	abcValues [3]int
}

func (o *opcodeComputer) tick() (done bool) {
	if o.registers[o.instructionPointer] >= len(o.instructions) {
		fmt.Println("Out of range instruction, terminating...")
		return true
	}
	instIndex := o.registers[o.instructionPointer]
	inst := o.instructions[instIndex]

	opcodeFunc := opcodeNamesToFuncs[inst.name]

	o.registers = opcodeFunc(o.registers, inst.abcValues)

	o.registers[o.instructionPointer]++

	if o.registers[o.instructionPointer] >= len(o.instructions) {
		return true
	}

	return false
}

func parseInput(input string) opcodeComputer {
	lines := strings.Split(input, "\n")

	var instructionPointer int
	fmt.Sscanf(lines[0], "#ip %d", &instructionPointer)

	var instructions []instruction
	for _, l := range lines[1:] {
		var inst instruction
		fmt.Sscanf(l, "%4s %d %d %d", &inst.name, &inst.abcValues[0], &inst.abcValues[1], &inst.abcValues[2])
		instructions = append(instructions, inst)
	}

	return opcodeComputer{
		instructions:       instructions,
		instructionPointer: instructionPointer,
	}
}

var opcodeNamesToFuncs = map[string]opcodeFunc{
	"addr": addr, "addi": addi,
	"mulr": mulr, "muli": muli,
	"banr": banr, "bani": bani,
	"borr": borr, "bori": bori,
	"setr": setr, "seti": seti,
	"gtir": gtir, "gtri": gtri, "gtrr": gtrr,
	"eqir": eqir, "eqri": eqri, "eqrr": eqrr,
}

type opcodeFunc func([6]int, [3]int) [6]int

func addr(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] + registers[abcValues[1]]
	return registers
}

func addi(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] + abcValues[1]
	return registers
}
func mulr(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] * registers[abcValues[1]]
	return registers
}
func muli(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] * abcValues[1]
	return registers
}
func banr(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] & registers[abcValues[1]]
	return registers
}
func bani(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] & abcValues[1]
	return registers
}
func borr(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] | registers[abcValues[1]]
	return registers
}
func bori(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]] | abcValues[1]
	return registers
}
func setr(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = registers[abcValues[0]]
	return registers
}
func seti(registers [6]int, abcValues [3]int) [6]int {
	registers[abcValues[2]] = abcValues[0]
	return registers
}
func gtir(registers [6]int, abcValues [3]int) [6]int {
	if abcValues[0] > registers[abcValues[1]] {
		registers[abcValues[2]] = 1
	} else {
		registers[abcValues[2]] = 0
	}
	return registers
}
func gtri(registers [6]int, abcValues [3]int) [6]int {
	if registers[abcValues[0]] > abcValues[1] {
		registers[abcValues[2]] = 1
	} else {
		registers[abcValues[2]] = 0
	}
	return registers
}
func gtrr(registers [6]int, abcValues [3]int) [6]int {
	if registers[abcValues[0]] > registers[abcValues[1]] {
		registers[abcValues[2]] = 1
	} else {
		registers[abcValues[2]] = 0
	}
	return registers
}
func eqir(registers [6]int, abcValues [3]int) [6]int {
	if abcValues[0] == registers[abcValues[1]] {
		registers[abcValues[2]] = 1
	} else {
		registers[abcValues[2]] = 0
	}
	return registers
}
func eqri(registers [6]int, abcValues [3]int) [6]int {
	if registers[abcValues[0]] == abcValues[1] {
		registers[abcValues[2]] = 1
	} else {
		registers[abcValues[2]] = 0
	}
	return registers
}
func eqrr(registers [6]int, abcValues [3]int) [6]int {
	if registers[abcValues[0]] == registers[abcValues[1]] {
		registers[abcValues[2]] = 1
	} else {
		registers[abcValues[2]] = 0
	}
	return registers
}
