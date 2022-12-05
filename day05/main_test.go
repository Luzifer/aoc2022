package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCargo(t *testing.T) {
	td, err := os.Open("./test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { td.Close() })

	c, err := readCargo(bufio.NewScanner(td))
	require.NoError(t, err, "reading cargo")

	assert.Equal(t, [][]byte{
		{'Z', 'N'},
		{'M', 'C', 'D'},
		{'P'},
	}, c.stacks, "expecting cargo to be read correctly")
}

func TestMoveOperation9000(t *testing.T) {
	td, err := os.Open("./test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { td.Close() })

	c, err := readCargo(bufio.NewScanner(td))
	require.NoError(t, err, "reading cargo")

	c.Move(2, 3, 2, moveMode9000)
	assert.Equal(t, [][]byte{
		{'Z', 'N'},
		{'M'},
		{'P', 'D', 'C'},
	}, c.stacks, "expecting cargo to have the move applied correctly")
}

func TestMoveOperation9001(t *testing.T) {
	td, err := os.Open("./test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { td.Close() })

	c, err := readCargo(bufio.NewScanner(td))
	require.NoError(t, err, "reading cargo")

	c.Move(2, 3, 2, moveMode9001)
	assert.Equal(t, [][]byte{
		{'Z', 'N'},
		{'M'},
		{'P', 'C', 'D'},
	}, c.stacks, "expecting cargo to have the move applied correctly")
}
