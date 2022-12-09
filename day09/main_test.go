package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestRope(hx, hy, tx, ty int) *rope {
	r := newRope(2)
	r.head.x = hx
	r.head.y = hy
	r.tail.x = tx
	r.tail.y = ty
	return r
}

func TestSegmentFollow(t *testing.T) {
	for _, tc := range []struct {
		r  *rope
		tt [2]int
	}{
		{newTestRope(3, 1, 1, 1), [2]int{2, 1}},
		{newTestRope(1, 1, 1, 3), [2]int{1, 2}},
		{newTestRope(2, 3, 1, 1), [2]int{2, 2}},
		{newTestRope(3, 2, 1, 1), [2]int{2, 2}},
	} {
		tc.r.applyRopePull()
		assert.Equal(t, tc.tt, tc.r.TailPos())
	}
}

func TestCoveredFields(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })

	r := newRope(2)
	r.ApplyMoveSet(f)

	assert.Equal(t, 13, r.CountCoveredFields(), "there are 13 positions the tail visited at least once")
}

func TestIsAdjacent(t *testing.T) {
	assert.True(t, newTestRope(0, 0, 0, 0).tail.isAdjacentToPrev(), "head == tail")

	assert.True(t, newTestRope(2, 1, 1, 1).tail.isAdjacentToPrev(), "tail left of head")
	assert.True(t, newTestRope(1, 1, 2, 1).tail.isAdjacentToPrev(), "tail right of head")
	assert.True(t, newTestRope(1, 2, 1, 1).tail.isAdjacentToPrev(), "tail below head")
	assert.True(t, newTestRope(1, 1, 1, 2).tail.isAdjacentToPrev(), "tail over head")

	assert.True(t, newTestRope(0, 0, 1, 1).tail.isAdjacentToPrev(), "tail right over head")
	assert.True(t, newTestRope(0, 0, -1, 1).tail.isAdjacentToPrev(), "tail left over head")
	assert.True(t, newTestRope(0, 0, -1, -1).tail.isAdjacentToPrev(), "tail left under head")
	assert.True(t, newTestRope(0, 0, 1, -1).tail.isAdjacentToPrev(), "tail right under head")
}

func TestMoveHead(t *testing.T) {
	r := newRope(2)
	assert.Equal(t, [2]int{0, 0}, r.HeadPos(), "expect rope to start head 0,0")
	assert.Equal(t, [2]int{0, 0}, r.TailPos(), "expect rope to start tail 0,0")

	r.MoveHead(directionRight, 4)
	assert.Equal(t, [2]int{4, 0}, r.HeadPos(), "expect rope head in defined position")
	assert.Equal(t, [2]int{3, 0}, r.TailPos(), "expect rope tail in defined position")

	r.MoveHead(directionUp, 4)
	assert.Equal(t, [2]int{4, 4}, r.HeadPos(), "expect rope head in defined position")
	assert.Equal(t, [2]int{4, 3}, r.TailPos(), "expect rope tail in defined position")

	r.MoveHead(directionLeft, 3)
	assert.Equal(t, [2]int{1, 4}, r.HeadPos(), "expect rope head in defined position")
	assert.Equal(t, [2]int{2, 4}, r.TailPos(), "expect rope tail in defined position")

	r.MoveHead(directionDown, 1)
	assert.Equal(t, [2]int{1, 3}, r.HeadPos(), "expect rope head in defined position")
	assert.Equal(t, [2]int{2, 4}, r.TailPos(), "expect rope tail in defined position")
}
