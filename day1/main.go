package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readInput() map[int]int64 {
	var (
		elfs    = map[int]int64{}
		elfNo   int
		scanner = bufio.NewScanner(os.Stdin)
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			elfNo++
			continue
		}

		cal, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			panic(err)
		}

		elfs[elfNo] += cal
	}

	return elfs
}

func main() {
	elfCalories := readInput()

	fmt.Printf("Solution 1: %d\n", solve1(elfCalories))
	fmt.Printf("Solution 2: %d\n", solve2(elfCalories))
}

func solve1(elfCalories map[int]int64) int64 {
	var maxCals int64

	for _, cal := range elfCalories {
		if cal > maxCals {
			maxCals = cal
		}
	}

	return maxCals
}

func solve2(elfCalories map[int]int64) int {
	var calories []int
	for _, cal := range elfCalories {
		calories = append(calories, int(cal))
	}

	sort.Ints(calories)

	top3 := calories[len(calories)-3:]
	var sum int
	for _, cal := range top3 {
		sum += cal
	}

	return sum
}
