package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var shapeTestBox = &shape{
	[]byte("#####"),
	[]byte("#   #"),
	[]byte("#   #"),
	[]byte("#   #"),
	[]byte("#####"),
}

func TestDoesCollide(t *testing.T) {
	c := newChamber()
	c.curShapePos = coord{0, 0}
	c.curShape = shapeTestBox
	c.width = 5
	c.settleShape()

	assert.True(t, shapePlus.DoesCollide(c, coord{0, 1}))
	assert.True(t, shapePlus.DoesCollide(c, coord{1, 0}))
	assert.True(t, shapePlus.DoesCollide(c, coord{1, 2}))
	assert.True(t, shapePlus.DoesCollide(c, coord{2, 1}))
	assert.False(t, shapePlus.DoesCollide(c, coord{1, 1}))
	assert.False(t, shapePlus.DoesCollide(c, coord{1, 5}))
}

func TestSettleShape(t *testing.T) {
	c := newChamber()
	c.width = 5
	c.curShapePos = coord{0, 0}
	c.curShape = shapeTestBox

	c.settleShape()

	assert.True(t, c.map2d[coord{0, 0}.String()])
	assert.True(t, c.map2d[coord{1, 0}.String()])
	assert.True(t, c.map2d[coord{2, 0}.String()])
	assert.True(t, c.map2d[coord{3, 0}.String()])
	assert.True(t, c.map2d[coord{4, 0}.String()])

	assert.True(t, c.map2d[coord{0, 1}.String()])
	assert.False(t, c.map2d[coord{1, 1}.String()])
	assert.False(t, c.map2d[coord{2, 1}.String()])
	assert.False(t, c.map2d[coord{3, 1}.String()])
	assert.True(t, c.map2d[coord{4, 1}.String()])
}

func TestSolve1(t *testing.T) {
	c := newChamber()
	c.width = 7
	c.jetDirs = jetDirsFromBytes([]byte(">>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>"))

	assert.Equal(t, 3068, heightAfter(c, 2022))
}

func TestTick(t *testing.T) {
	c := newChamber()
	c.width = 7
	c.jetDirs = jetDirsFromBytes([]byte(">>><<><>><<<>><>>><<<>>><<<><<<>><>><<>>"))

	for i := 0; i < 100; i++ {
		c.Tick()
	}

	buf := new(bytes.Buffer)
	c.Draw(buf)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")

	exp := strings.TrimSpace(`
|##··##·|
|######·|
|·###···|
|··#····|
|·####··|
|····##·|
|····##·|
|····#··|
|··#·#··|
|··#·#··|
|#####··|
|··###··|
|···#···|
|··####·|
`)

	require.GreaterOrEqual(t, len(lines), 14, "output should be taller than 14l")
	assert.Equal(t, exp, strings.Join(lines[len(lines)-14:], "\n"))
}
