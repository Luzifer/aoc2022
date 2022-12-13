package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testReadSegment(line string) segment {
	s, _ := readSegment(line)
	return s
}

func TestCorrectPacketIndexSum(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	call := readDistressCall(f)

	assert.Equal(t, 13, call.CorrectPacketIndexSum())
}

func TestReadDistressCall(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	call := readDistressCall(f)

	assert.Equal(t, call, distressCall{
		{readPacket("[1,1,3,1,1]"), readPacket("[1,1,5,1,1]")},
		{readPacket("[[1],[2,3,4]]"), readPacket("[[1],4]")},
		{readPacket("[9]"), readPacket("[[8,7,6]]")},
		{readPacket("[[4,4],4,4]"), readPacket("[[4,4],4,4,4]")},
		{readPacket("[7,7,7,7]"), readPacket("[7,7,7]")},
		{readPacket("[]"), readPacket("[3]")},
		{readPacket("[[[]]]"), readPacket("[[]]")},
		{readPacket("[1,[2,[3,[4,[5,6,7]]]],8,9]"), readPacket("[1,[2,[3,[4,[5,6,0]]]],8,9]")},
	})
}

func TestReadSegment(t *testing.T) {
	for line, exp := range map[string]listSegment{
		"[9]": {
			intSegment(9),
		},
		"[1,1,3,1,1]": {
			intSegment(1),
			intSegment(1),
			intSegment(3),
			intSegment(1),
			intSegment(1),
		},
		"[[1],[2,3,4]]": {
			listSegment{intSegment(1)},
			listSegment{
				intSegment(2),
				intSegment(3),
				intSegment(4),
			},
		},
		"[[[]]]": {listSegment{listSegment{}}},
		"[[4,4],4,4,4]": {
			listSegment{intSegment(4), intSegment(4)},
			intSegment(4),
			intSegment(4),
			intSegment(4),
		},
	} {
		seg, n := readSegment(line)
		assert.Equal(t, len(line), n, "expect line to be consumed")
		assert.Equal(t, exp, seg, "expect segment to be correct")
	}
}

func TestSegmentCompareTo(t *testing.T) {
	for _, tc := range []struct {
		s1, s2 segment
		result compareResult
	}{
		{intSegment(1), intSegment(2), compareResultSmaller},
		{intSegment(2), intSegment(1), compareResultBigger},
		{intSegment(2), intSegment(2), compareResultEqual},
		{listSegment{intSegment(2)}, listSegment{intSegment(2)}, compareResultEqual},
		{listSegment{intSegment(2), intSegment(1)}, listSegment{intSegment(2)}, compareResultBigger},
		{listSegment{intSegment(2)}, listSegment{intSegment(2), intSegment(1)}, compareResultSmaller},

		{testReadSegment("[1,1,3,1,1]"), testReadSegment("[1,1,5,1,1]"), compareResultSmaller},
		{testReadSegment("[[1],[2,3,4]]"), testReadSegment("[[1],4]"), compareResultSmaller},
		{testReadSegment("[9]"), testReadSegment("[[8,7,6]]"), compareResultBigger},
		{testReadSegment("[[4,4],4,4]"), testReadSegment("[[4,4],4,4,4]"), compareResultSmaller},
		{testReadSegment("[7,7,7,7]"), testReadSegment("[7,7,7]"), compareResultBigger},
		{testReadSegment("[]"), testReadSegment("[3]"), compareResultSmaller},
		{testReadSegment("[[[]]]"), testReadSegment("[[]]"), compareResultBigger},
		{testReadSegment("[1,[2,[3,[4,[5,6,7]]]],8,9]"), testReadSegment("[1,[2,[3,[4,[5,6,0]]]],8,9]"), compareResultBigger},
	} {
		assert.Equal(t, tc.result, tc.s1.CompareTo(tc.s2))
	}
}

func TestSolve2(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	call := readDistressCall(f)

	assert.Equal(t, 140, solve2(call))
}
