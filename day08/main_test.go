package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadGrid(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	g := readGrid(f)

	assert.Equal(t, 4, g.MaxX(), "testing grid width")
	assert.Equal(t, 4, g.MaxY(), "testing grid height")

	assert.Equal(t, tree{Height: 3, X: 3, Y: 2}, g.GetTree(3, 2), "testing tree sample")
}

func TestGetScenicScore(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	g := readGrid(f)

	assert.Equal(t, 4, g.GetTree(2, 1).GetScenicScore(g), "For this tree, this is 4")
	assert.Equal(t, 8, g.GetTree(2, 3).GetScenicScore(g), "This tree's scenic score is 8")
}

func TestIsOutsideVisible(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	g := readGrid(f)

	assert.True(t, g.GetTree(1, 1).IsOutsideVisible(g), "top-left 5 is visible")
	assert.True(t, g.GetTree(2, 1).IsOutsideVisible(g), "top-middle 5 is visible")
	assert.False(t, g.GetTree(3, 1).IsOutsideVisible(g), "top-right 1 is not visible")
	assert.True(t, g.GetTree(1, 2).IsOutsideVisible(g), "left-middle 5 is visible")
	assert.False(t, g.GetTree(2, 2).IsOutsideVisible(g), "center 3 is not visible")
	assert.True(t, g.GetTree(3, 2).IsOutsideVisible(g), "right-middle 3 is visible")
	assert.False(t, g.GetTree(1, 3).IsOutsideVisible(g), "but the 3 and 4 are not")
	assert.True(t, g.GetTree(2, 3).IsOutsideVisible(g), "bottom row, the middle 5 is visible")
	assert.False(t, g.GetTree(3, 3).IsOutsideVisible(g), "but the 3 and 4 are not")
}
