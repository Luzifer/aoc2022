package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type (
	direction byte

	rope struct {
		head, tail        *segment
		tailCoveredCoords map[string]bool
	}

	segment struct {
		x, y       int
		next, prev *segment
	}
)

const (
	directionDown  direction = 'D'
	directionLeft  direction = 'L'
	directionRight direction = 'R'
	directionUp    direction = 'U'
)

func newRope(l int) *rope {
	var head, last, tail *segment
	for i := 0; i < l; i++ {
		seg := &segment{}
		if head == nil {
			head = seg
		}
		if last != nil {
			last.next = seg
			seg.prev = last
		}
		last = seg
		tail = seg
	}

	r := &rope{head: head, tail: tail}
	r.recordTail()
	return r
}

func (r *rope) ApplyMoveSet(i io.Reader) {
	scanner := bufio.NewScanner(i)
	for scanner.Scan() {
		dir := direction(scanner.Text()[0])
		steps, _ := strconv.Atoi(scanner.Text()[2:])

		r.MoveHead(dir, steps)
	}
}

func (r rope) CountCoveredFields() int { return len(r.tailCoveredCoords) }

func (r rope) HeadPos() [2]int { return [2]int{r.head.x, r.head.y} }

func (r *rope) MoveHead(dir direction, steps int) {
	for i := 0; i < steps; i++ {
		switch dir {
		case directionDown:
			r.head.y -= 1
		case directionLeft:
			r.head.x -= 1
		case directionRight:
			r.head.x += 1
		case directionUp:
			r.head.y += 1
		}

		r.applyRopePull()
	}
}

func (r rope) TailPos() [2]int { return [2]int{r.tail.x, r.tail.y} }

func (r *rope) applyRopePull() {
	seg := r.head
	for seg != nil {
		seg.follow()
		seg = seg.next
	}

	r.recordTail()
}

func (r *rope) recordTail() {
	if r.tailCoveredCoords == nil {
		r.tailCoveredCoords = make(map[string]bool)
	}
	r.tailCoveredCoords[fmt.Sprintf("%d:%d", r.tail.x, r.tail.y)] = true
}

func (s *segment) follow() {
	if s.prev == nil {
		// Nothing to follow
		return
	}

	if s.isAdjacentToPrev() {
		// No action needed
		return
	}

	switch {
	case s.x == s.prev.x:
		// X axis is the same, we need to move tail up / down
		if s.y < s.prev.y {
			s.y += 1
		} else {
			s.y -= 1
		}

	case s.y == s.prev.y:
		// Y axis is the same, we need to move tail left / right
		if s.x < s.prev.x {
			s.x += 1
		} else {
			s.x -= 1
		}

	default:
		// Both axis differ, we need to make a diagonal step
		if s.y < s.prev.y {
			s.y += 1
		} else {
			s.y -= 1
		}

		if s.x < s.prev.x {
			s.x += 1
		} else {
			s.x -= 1
		}
	}
}

func (s segment) isAdjacentToPrev() bool {
	return int(math.Sqrt(math.Pow(float64(s.x-s.prev.x), 2)+math.Pow(float64(s.y-s.prev.y), 2))) < 2
}

func main() {
	moveSet, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	r := newRope(2)
	r.ApplyMoveSet(bytes.NewReader(moveSet))
	fmt.Printf("Solution 1: %d\n", r.CountCoveredFields())

	r = newRope(10)
	r.ApplyMoveSet(bytes.NewReader(moveSet))
	fmt.Printf("Solution 2: %d\n", r.CountCoveredFields())
}
