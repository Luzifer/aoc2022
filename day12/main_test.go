package main

import (
	"log"
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

	assert.Equal(t, int('a'), g.ElevationAt(0, 0))
	assert.Equal(t, int('q'), g.ElevationAt(3, 0))
	assert.Equal(t, int('u'), g.ElevationAt(4, 3))
	assert.Equal(t, int('z'), g.ElevationAt(5, 2))
	assert.Equal(t, int('i'), g.ElevationAt(7, 4))
}

func TestShortestPath(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	g := readGrid(f)

	p, _, _ := g.ShortestPath(g.TileAt(g.start[0], g.start[1]))
	for _, pe := range p {
		log.Printf("%#v", pe)
	}
}

func TestShortestPathLen(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	g := readGrid(f)

	pl, _ := g.ShortestPathLen(g.TileAt(g.start[0], g.start[1]))
	assert.Equal(t, 31, pl, "This path reaches the goal in 31 steps, the fewest possible.")
}

func TestFindShortestACoordPathLen(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	g := readGrid(f)

	start, pl := g.FindShortestACoordPathLen()
	p, _, _ := g.ShortestPath(start)
	for _, pe := range p {
		log.Printf("%#v", pe)
	}

	assert.Equal(t, 29, pl, "This path reaches the goal in only 29 steps, the fewest possible.")
}
