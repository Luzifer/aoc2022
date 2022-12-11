package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type (
	monkey struct {
		Number      uint64
		Items       []uint64
		Operation   func(itemLevel uint64) uint64
		TestDivisor uint64
		TargetTrue  uint64
		TargetFalse uint64

		itemInspected uint64
	}

	monkeyGroup []*monkey
)

const (
	inputRegexpGroupMonkeyNo = iota + 1
	inputRegexpGroupStartingItems
	inputRegexpGroupOpNumber1
	inputRegexpGroupOpAction
	inputRegexpGroupOpNumber2
	inputRegexpGroupTestDivisor
	inputRegexpGroupTargetTrue
	inputRegexpGroupTargetFalse
)

var inputRegexp = regexp.MustCompile(`Monkey ([0-9]+):\s*Starting items: (?:([0-9, ]+))\s*Operation: new = (old|[0-9]+) ([/*+-]) (old|[0-9]+)\s*Test: divisible by ([0-9]+)\s*If true: throw to monkey ([0-9]+)\s*If false: throw to monkey ([0-9]+)`)

func readMonkeyGroup(r io.Reader) monkeyGroup {
	textform, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	var (
		matches = inputRegexp.FindAllStringSubmatch(string(textform), -1)
		mg      monkeyGroup
	)
	for _, match := range matches {
		m := &monkey{}
		m.Number, _ = strconv.ParseUint(match[inputRegexpGroupMonkeyNo], 10, 64)
		m.TestDivisor, _ = strconv.ParseUint(match[inputRegexpGroupTestDivisor], 10, 64)
		m.TargetTrue, _ = strconv.ParseUint(match[inputRegexpGroupTargetTrue], 10, 64)
		m.TargetFalse, _ = strconv.ParseUint(match[inputRegexpGroupTargetFalse], 10, 64)

		for _, il := range strings.Split(match[inputRegexpGroupStartingItems], ", ") {
			l, _ := strconv.ParseUint(il, 10, 64)
			m.Items = append(m.Items, l)
		}

		if match[inputRegexpGroupOpNumber2] == "old" {
			switch match[inputRegexpGroupOpAction] {
			case "+":
				m.Operation = func(il uint64) uint64 { return il + il }
			case "-":
				m.Operation = func(il uint64) uint64 { return il - il }
			case "*":
				m.Operation = func(il uint64) uint64 { return il * il }
			case "/":
				m.Operation = func(il uint64) uint64 { return il / il }
			}
		} else {
			n2, _ := strconv.ParseUint(match[inputRegexpGroupOpNumber2], 10, 64)
			switch match[inputRegexpGroupOpAction] {
			case "+":
				m.Operation = func(il uint64) uint64 { return il + n2 }
			case "-":
				m.Operation = func(il uint64) uint64 { return il - n2 }
			case "*":
				m.Operation = func(il uint64) uint64 { return il * n2 }
			case "/":
				m.Operation = func(il uint64) uint64 { return il / n2 }
			}
		}

		mg = append(mg, m)
	}

	return mg
}

func (m *monkey) Catch(il uint64) { m.Items = append(m.Items, il) }

func (m *monkey) ExecuteTurn(mg monkeyGroup, worryDivisor uint64) {
	modifier := func(in uint64) uint64 { return in / worryDivisor }

	if worryDivisor == 0 {
		// Used to reduce the frickin large numbers (> max-uint64) into
		// "managable" numbers. If a number was dividable through all
		// of the divisors before it still is when taken the modulo
		// of the product of the divisors.
		//
		// Knowledge taken from the AoC2022 solution and code ported
		// from a Rust solution as I'm lacking the mathematical knowledge
		// for this.

		var divisor uint64 = 1
		for _, m := range mg {
			divisor *= m.TestDivisor
		}
		modifier = func(in uint64) uint64 { return in % divisor }
	}

	for _, il := range m.Items {
		m.itemInspected++

		il = m.Operation(il)
		il = modifier(il)

		if il%m.TestDivisor == 0 {
			mg[m.TargetTrue].Catch(il)
		} else {
			mg[m.TargetFalse].Catch(il)
		}
	}

	m.Items = nil
}

func (m monkey) InspectedItems() uint64 { return m.itemInspected }

func (m monkeyGroup) ExecuteRounds(n int, worryDivisor uint64) {
	for r := 0; r < n; r++ {
		for i := range m {
			m[i].ExecuteTurn(m, worryDivisor)
		}
	}
}

func main() {
	buf := new(bytes.Buffer)
	io.Copy(buf, os.Stdin)

	mg := readMonkeyGroup(bytes.NewReader(buf.Bytes()))
	mg.ExecuteRounds(20, 3)

	mg2 := readMonkeyGroup(bytes.NewReader(buf.Bytes()))
	mg2.ExecuteRounds(10000, 0)

	fmt.Printf("Solution 1: %d\n", solutionFromGroup(mg))
	fmt.Printf("Solution 2: %d\n", solutionFromGroup(mg2))
}

func solutionFromGroup(mg monkeyGroup) uint64 {
	inspectCounts := []uint64{}
	for _, m := range mg {
		inspectCounts = append(inspectCounts, m.InspectedItems())
	}
	sort.SliceStable(inspectCounts, func(i, j int) bool { return inspectCounts[i] > inspectCounts[j] })

	var solution uint64 = 1
	for _, ic := range inspectCounts[:2] {
		solution *= ic
	}

	return solution
}
