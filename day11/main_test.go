package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadMonkeyGroup(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test input")
	t.Cleanup(func() { f.Close() })

	mg := readMonkeyGroup(f)
	assert.Equal(t, []uint64{79, 98}, mg[0].Items)
	assert.Equal(t, uint64(2*19), mg[0].Operation(2))
	assert.Equal(t, uint64(23), mg[0].TestDivisor)
	assert.Equal(t, uint64(2), mg[0].TargetTrue)
	assert.Equal(t, uint64(3), mg[0].TargetFalse)

	assert.Equal(t, []uint64{54, 65, 75, 74}, mg[1].Items)
	assert.Equal(t, uint64(2+6), mg[1].Operation(2))
	assert.Equal(t, uint64(19), mg[1].TestDivisor)
	assert.Equal(t, uint64(2), mg[1].TargetTrue)
	assert.Equal(t, uint64(0), mg[1].TargetFalse)

	assert.Equal(t, []uint64{79, 60, 97}, mg[2].Items)
	assert.Equal(t, uint64(2*2), mg[2].Operation(2))
	assert.Equal(t, uint64(13), mg[2].TestDivisor)
	assert.Equal(t, uint64(1), mg[2].TargetTrue)
	assert.Equal(t, uint64(3), mg[2].TargetFalse)

	assert.Equal(t, []uint64{74}, mg[3].Items)
	assert.Equal(t, uint64(2+3), mg[3].Operation(2))
	assert.Equal(t, uint64(17), mg[3].TestDivisor)
	assert.Equal(t, uint64(0), mg[3].TargetTrue)
	assert.Equal(t, uint64(1), mg[3].TargetFalse)
}

func TestRoundExecution(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test input")
	t.Cleanup(func() { f.Close() })

	mg := readMonkeyGroup(f)
	mg.ExecuteRounds(1, 3)

	assert.Equal(t, []uint64{20, 23, 27, 26}, mg[0].Items)
	assert.Equal(t, []uint64{2080, 25, 167, 207, 401, 1046}, mg[1].Items)
	assert.Equal(t, ([]uint64)(nil), mg[2].Items)
	assert.Equal(t, ([]uint64)(nil), mg[3].Items)

	mg.ExecuteRounds(9, 3) // We're now "after round 10"
	assert.Equal(t, []uint64{91, 16, 20, 98}, mg[0].Items)
	assert.Equal(t, []uint64{481, 245, 22, 26, 1092, 30}, mg[1].Items)
	assert.Equal(t, ([]uint64)(nil), mg[2].Items)
	assert.Equal(t, ([]uint64)(nil), mg[3].Items)

	mg.ExecuteRounds(10, 3) // We're now "after round 20"
	assert.Equal(t, []uint64{10, 12, 14, 26, 34}, mg[0].Items)
	assert.Equal(t, []uint64{245, 93, 53, 199, 115}, mg[1].Items)
	assert.Equal(t, ([]uint64)(nil), mg[2].Items)
	assert.Equal(t, ([]uint64)(nil), mg[3].Items)

	// Count the total number of times each monkey inspects items over 20 rounds:
	assert.Equal(t, uint64(101), mg[0].itemInspected)
	assert.Equal(t, uint64(95), mg[1].itemInspected)
	assert.Equal(t, uint64(7), mg[2].itemInspected)
	assert.Equal(t, uint64(105), mg[3].itemInspected)
}

func TestSolution2(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test input")
	t.Cleanup(func() { f.Close() })

	mg := readMonkeyGroup(f)
	mg.ExecuteRounds(1, 0)
	assert.Equal(t, []uint64{2, 4, 3, 6}, []uint64{mg[0].itemInspected, mg[1].itemInspected, mg[2].itemInspected, mg[3].itemInspected}, "round 1")

	mg.ExecuteRounds(19, 0)
	assert.Equal(t, []uint64{99, 97, 8, 103}, []uint64{mg[0].itemInspected, mg[1].itemInspected, mg[2].itemInspected, mg[3].itemInspected}, "round 20")

	mg.ExecuteRounds(980, 0)
	assert.Equal(t, []uint64{5204, 4792, 199, 5192}, []uint64{mg[0].itemInspected, mg[1].itemInspected, mg[2].itemInspected, mg[3].itemInspected}, "round 1000")

	mg.ExecuteRounds(9000, 0)
	assert.Equal(t, uint64(2713310158), solutionFromGroup(mg))
}
