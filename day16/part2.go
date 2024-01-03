package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("day16/input.txt")
	if err != nil {
		panic(err)
	}
	inputStr := strings.TrimSpace(string(input))
	lines := strings.Split(inputStr, "\n")

	opcodes := []OP{
		{name: "addr", action: '+', a: 'r', b: 'r'},
		{name: "addi", action: '+', a: 'r', b: 'v'},
		{name: "mulr", action: '*', a: 'r', b: 'r'},
		{name: "muli", action: '*', a: 'r', b: 'v'},
		{name: "banr", action: '&', a: 'r', b: 'r'},
		{name: "bani", action: '&', a: 'r', b: 'v'},
		{name: "borr", action: '|', a: 'r', b: 'r'},
		{name: "bori", action: '|', a: 'r', b: 'v'},
		{name: "setr", action: 'a', a: 'r', b: 'r'},
		{name: "seti", action: 'a', a: 'v', b: 'r'},
		{name: "gtir", action: '>', a: 'v', b: 'r'},
		{name: "gtri", action: '>', a: 'r', b: 'v'},
		{name: "gtrr", action: '>', a: 'r', b: 'r'},
		{name: "eqir", action: '=', a: 'v', b: 'r'},
		{name: "eqri", action: '=', a: 'r', b: 'v'},
		{name: "eqir", action: '=', a: 'r', b: 'r'},
	}

	sum := 0
	var lineCount int
	for lineCount < len(lines) {
		if len(lines[lineCount]) > 0 && lines[lineCount][0] == 'B' {
			split := regSplit(lines[lineCount], "[^0-9]+")
			registers := []int{
				strToInt(split[1]),
				strToInt(split[2]),
				strToInt(split[3]),
				strToInt(split[4]),
			}
			split = regSplit(lines[lineCount+1], "[^0-9]+")
			instruction := []byte{
				byte(strToInt(split[0])),
				byte(strToInt(split[1])),
				byte(strToInt(split[2])),
				byte(strToInt(split[3])),
			}
			split = regSplit(lines[lineCount+2], "[^0-9]+")
			result := []int{
				strToInt(split[1]),
				strToInt(split[2]),
				strToInt(split[3]),
				strToInt(split[4]),
			}
			tempSum := testCode(registers, result, instruction, opcodes)

			if tempSum >= 3 {
				sum++
			}

			lineCount = lineCount + 4
		} else {
			break
		}
	}

	orderedOpCodes := make(map[byte]*OP, 0)

	for len(orderedOpCodes) < 16 {
		for i := range opcodes {

			if len(opcodes[i].matchCount) == 1 {
				c := opcodes[i].matchCount[0]
				orderedOpCodes[c] = &opcodes[i]
				for j := range opcodes {
					remove(&opcodes[j], c)
				}
			}
		}
	}

	lineCount = lineCount + 2

	r := make([]int, 4)

	for ; lineCount < len(lines); lineCount++ {
		split := regSplit(lines[lineCount], "[^0-9]+")
		instruction := []byte{
			byte(strToInt(split[0])),
			byte(strToInt(split[1])),
			byte(strToInt(split[2])),
			byte(strToInt(split[3])),
		}

		r = runOp(*orderedOpCodes[instruction[0]], r, instruction)
	}

	fmt.Println(r[0])
}

func remove(op *OP, c byte) {
	i := -1
	for j, v := range op.matchCount {
		if c == v {
			i = j
		}
	}
	if i != -1 {
		op.matchCount = append(op.matchCount[:i], op.matchCount[i+1:]...)
	}
}

func add(op *OP, c byte) {
	for _, v := range op.matchCount {
		if v == c {
			return
		}
	}
	op.matchCount = append(op.matchCount, c)
}

func testCode(registers, result []int, instruction []byte, opcodes []OP) int {
	sum := 0
	for i := range opcodes {
		if match(result, runOp(opcodes[i], registers, instruction)) {
			add(&opcodes[i], instruction[0])
			sum++
		}
	}
	return sum
}

func match(r, c []int) bool {
	if len(r) != len(c) {
		return false
	}
	for i := range r {
		if r[i] != c[i] {
			return false
		}
	}
	return true
}

func runOp(op OP, registers []int, instruction []byte) []int {
	registerCP := make([]int, 4)
	copy(registerCP, registers)
	var A, B int
	if op.a == 'r' {
		A = registerCP[instruction[1]]
	} else {
		A = int(instruction[1])
	}
	if op.b == 'r' {
		B = registerCP[instruction[2]]
	} else {
		B = int(instruction[2])
	}
	switch op.action {
	case '+':
		registerCP[instruction[3]] = A + B
		break
	case '*':
		registerCP[instruction[3]] = A * B
		break
	case '&':
		registerCP[instruction[3]] = A & B
		break
	case '|':
		registerCP[instruction[3]] = A | B
		break
	case 'a':
		registerCP[instruction[3]] = A
		break
	case '>':
		if A > B {
			registerCP[instruction[3]] = 1
		} else {
			registerCP[instruction[3]] = 0
		}
		break
	case '=':
		if A == B {
			registerCP[instruction[3]] = 1
		} else {
			registerCP[instruction[3]] = 0
		}
		break
	default:
		fmt.Println("not valid instruction")
	}
	return registerCP
}

type OP struct {
	a          rune
	b          rune
	action     rune
	name       string
	matchCount []byte
}

func strToInt(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func regSplit(text string, delimeter string) []string {
	reg := regexp.MustCompile(delimeter)
	i := reg.FindAllStringIndex(text, -1)
	laststart := 0
	n := make([]string, len(i)+1)
	for i, element := range i {
		n[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	n[len(i)] = text[laststart:]
	return n
}
