package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type (
	distressCall  [][]packet
	compareResult uint

	intSegment  int
	listSegment []segment

	packet listSegment

	segment interface {
		CompareTo(segment) compareResult
		IsList() bool
	}
)

const (
	compareResultEqual compareResult = iota
	compareResultSmaller
	compareResultBigger
)

var (
	_ segment = intSegment(0)
	_ segment = listSegment{}
)

func readDistressCall(r io.Reader) distressCall {
	var (
		call    distressCall
		pair    []packet
		scanner = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		if scanner.Text() == "" {
			call = append(call, pair)
			pair = nil
			continue
		}

		pair = append(pair, readPacket(scanner.Text()))
	}

	call = append(call, pair)

	return call
}

func (d distressCall) CorrectPacketIndexSum() int {
	var sum int

	for i := range d {
		if d[i][0].CompareTo(d[i][1]) == compareResultSmaller {
			sum += i + 1
		}
	}

	return sum
}

func (d distressCall) PacketList() []packet {
	var pl []packet
	for _, pair := range d {
		pl = append(pl, pair...)
	}

	return pl
}

func readPacket(line string) packet {
	s, n := readSegment(line)
	if n != len(line) {
		panic("packet line was not fully consumed by readSegment")
	}

	return packet(s.(listSegment))
}

func (p packet) CompareTo(t packet) compareResult {
	return listSegment(p).CompareTo(listSegment(t))
}

func readSegment(line string) (segment, int) {
	if line[0] == ']' {
		// We're in an empty list
		return nil, 0
	}

	if line[0] != '[' {
		// Int-Segment
		commaPos := float64(strings.IndexByte(line, ','))
		if commaPos < 0 {
			commaPos = math.MaxFloat64
		}
		rbPos := float64(strings.IndexByte(line, ']'))
		if commaPos < 0 {
			rbPos = math.MaxFloat64
		}
		charPos := math.Min(commaPos, rbPos)
		readUntil := math.Min(charPos, float64(len(line)-1))

		// log.Printf("reading int from %q: comma=%.0f rb=%.0f char=%.0f ru=%.0f", line, commaPos, rbPos, charPos, readUntil)

		number := line[:int(readUntil)]
		i, _ := strconv.Atoi(number)
		return intSegment(i), len(number)
	}

	// List-Segment
	var (
		lp int = 1
		ls     = listSegment{}
	)
	for {
		s, n := readSegment(line[lp:])
		lp += n
		if s != nil {
			ls = append(ls, s)
		}

		if line[lp] == ',' {
			lp += 1
		}

		if line[lp] == ']' {
			return ls, lp + 1
		}
	}
}

func (s intSegment) CompareTo(t segment) compareResult {
	if !t.IsList() {
		// Simple numeric comparison
		switch {
		case int(s) == int(t.(intSegment)):
			return compareResultEqual
		case int(s) < int(t.(intSegment)):
			return compareResultSmaller
		case int(s) > int(t.(intSegment)):
			return compareResultBigger
		}
		panic("unreachable")
	}

	// Other one is a list, we need to be one too
	return listSegment{s}.CompareTo(t)
}

func (s intSegment) IsList() bool { return false }

func (s listSegment) CompareTo(t segment) compareResult {
	if !t.IsList() {
		// Other one is a number, we are not, use them as list
		return s.CompareTo(listSegment{t})
	}

	// Both are (now) lists!
	tL := t.(listSegment)
	sN := len(s)
	tN := len(tL)
	cN := int(math.Min(float64(sN), float64(tN)))
	for i := 0; i < cN; i++ {
		if r := s[i].CompareTo(tL[i]); r != compareResultEqual {
			return r
		}
	}

	// All comparable items were equal, let list length decide
	switch {
	case sN == tN:
		// Lists have equal length
		return compareResultEqual
	case sN < tN:
		return compareResultSmaller
	case sN > tN:
		return compareResultBigger
	}

	panic("unreachable2")
}

func (s listSegment) IsList() bool { return true }

func main() {
	call := readDistressCall(os.Stdin)
	fmt.Printf("Solution 1: %d\n", call.CorrectPacketIndexSum())
	fmt.Printf("Solution 2: %d\n", solve2(call))
}

func solve2(call distressCall) int {
	marker1 := readPacket("[[2]]")
	marker2 := readPacket("[[6]]")

	pl := call.PacketList()
	pl = append(pl, marker1, marker2)
	sort.Slice(pl, func(i, j int) bool { return pl[i].CompareTo(pl[j]) == compareResultSmaller })

	var solution int = 1
	for i := range pl {
		if reflect.DeepEqual(marker1, pl[i]) || reflect.DeepEqual(marker2, pl[i]) {
			solution *= i + 1
		}
	}

	return solution
}
