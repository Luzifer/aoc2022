package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadInput(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test input")
	t.Cleanup(func() { f.Close() })

	g := newGrid()
	g.ReadInput(f)

	assert.Equal(t, fieldTypeRock, g.MaterialAt(503, 4))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(502, 4))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(502, 5))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(502, 6))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(502, 7))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(502, 8))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(502, 9))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(501, 9))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(500, 9))
	assert.Equal(t, fieldTypeRock, g.MaterialAt(494, 9))

	assert.Equal(t, fieldTypeAir, g.MaterialAt(504, 4))
	assert.Equal(t, fieldTypeAir, g.MaterialAt(493, 9))

	g.Render(os.Stdout)
}

func TestSettleSandGrains(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test input")
	t.Cleanup(func() { f.Close() })

	g := newGrid()
	g.ReadInput(f)

	n, err := g.SettleSandGrains(1, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, fieldTypeAir, g.MaterialAt(499, 8))
	assert.Equal(t, fieldTypeSand, g.MaterialAt(500, 8))
	assert.Equal(t, fieldTypeAir, g.MaterialAt(501, 8))

	n, err = g.SettleSandGrains(1, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, fieldTypeSand, g.MaterialAt(499, 8))
	assert.Equal(t, fieldTypeSand, g.MaterialAt(500, 8))
	assert.Equal(t, fieldTypeAir, g.MaterialAt(501, 8))

	n, err = g.SettleSandGrains(1, false)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.Equal(t, fieldTypeSand, g.MaterialAt(499, 8))
	assert.Equal(t, fieldTypeSand, g.MaterialAt(500, 8))
	assert.Equal(t, fieldTypeSand, g.MaterialAt(501, 8))

	n, err = g.SettleSandGrains(500, false)
	assert.ErrorIs(t, err, errSandFellToVoid)
	assert.Equal(t, 21, n)
}

func TestSettleSandGrainsWithFloor(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test input")
	t.Cleanup(func() { f.Close() })

	g := newGrid()
	g.ReadInput(f)

	n, err := g.SettleSandGrains(10000, true)
	assert.ErrorIs(t, err, errSandInputClogged)
	assert.Equal(t, 93, n)
}
