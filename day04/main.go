package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	elfPair struct {
		elf1 [2]int
		elf2 [2]int
	}
)

func readElfPair(line string) (out elfPair) {
	elfs := strings.Split(line, ",")

	elf1Range := strings.Split(elfs[0], "-")
	out.elf1[0], _ = strconv.Atoi(elf1Range[0])
	out.elf1[1], _ = strconv.Atoi(elf1Range[1])

	elf2Range := strings.Split(elfs[1], "-")
	out.elf2[0], _ = strconv.Atoi(elf2Range[0])
	out.elf2[1], _ = strconv.Atoi(elf2Range[1])

	return out
}

func (e elfPair) HasFullOverlap() bool {
	return (e.elf1[0] <= e.elf2[0] && e.elf2[1] <= e.elf1[1]) ||
		(e.elf2[0] <= e.elf1[0] && e.elf1[1] <= e.elf2[1])
}

func (e elfPair) HasOverlap() bool {
	return e.between(e.elf2[0], e.elf1[0], e.elf1[1]) ||
		e.between(e.elf2[1], e.elf1[0], e.elf1[1]) ||
		e.between(e.elf1[0], e.elf2[0], e.elf2[1]) ||
		e.between(e.elf1[1], e.elf2[0], e.elf2[1])
}

func (e elfPair) between(search, rangeL, rangeR int) bool {
	return rangeL <= search && search <= rangeR
}

func main() {
	var solution1, solution2 int

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		pair := readElfPair(scanner.Text())

		if pair.HasFullOverlap() {
			solution1++
		}

		if pair.HasOverlap() {
			solution2++
		}
	}

	fmt.Printf("Solution 1: %d\n", solution1)
	fmt.Printf("Solution 2: %d\n", solution2)
}
