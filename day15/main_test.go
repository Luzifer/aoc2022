package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadInput(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })

	a := readSensorArray(f)

	assert.Equal(t, sensorArray{
		&sensor{coord{2, 18}, coord{-2, 15}},
		&sensor{coord{9, 16}, coord{10, 16}},
		&sensor{coord{13, 2}, coord{15, 3}},
		&sensor{coord{12, 14}, coord{10, 16}},
		&sensor{coord{10, 20}, coord{10, 16}},
	}, a[:5])
}

func TestSolve1(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })

	a := readSensorArray(f)

	assert.Equal(t, 26, solve1(a, 10))
}

func TestSolve2(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err)
	t.Cleanup(func() { f.Close() })

	a := readSensorArray(f)

	assert.Equal(t, int64(56000011), solve2(a, 20))
}
