package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/beefsack/go-astar"
)

type (
	grid struct {
		terrain       [][]*tile
		start, target [2]int
	}

	tile struct {
		elevation int
		grid      *grid
		x, y      int
	}
)

var _ astar.Pather = tile{}

func readGrid(r io.Reader) *grid {
	var (
		out = &grid{}
		y   int
	)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		out.terrain = append(out.terrain, nil)
		for i, t := range scanner.Bytes() {
			switch t {
			case 'S':
				out.start = [2]int{i, y}
				t = 'a'

			case 'E':
				out.target = [2]int{i, y}
				t = 'z'
			}

			out.terrain[y] = append(out.terrain[y], &tile{elevation: int(t), grid: out, x: i, y: y})
		}

		y++
	}

	return out
}

func (g grid) ElevationAt(x, y int) int { return g.TileAt(x, y).elevation }
func (g grid) MaxX() int                { return len(g.terrain[0]) - 1 }
func (g grid) MaxY() int                { return len(g.terrain) - 1 }
func (g grid) TileAt(x, y int) *tile    { return g.terrain[y][x] }

func (g grid) FindShortestACoordPathLen() (*tile, int) {
	var (
		shortest      = math.MaxInt
		shortestStart *tile
	)
	for y := 0; y <= g.MaxY(); y++ {
		for x := 0; x <= g.MaxX(); x++ {
			possibleStart := g.TileAt(x, y)
			if possibleStart.elevation != int('a') {
				// That's not a start point
				continue
			}

			if l, ok := g.ShortestPathLen(possibleStart); ok && l < shortest {
				shortest = l
				shortestStart = possibleStart
			}
		}
	}

	return shortestStart, shortest
}

func (g grid) Print(w io.Writer, renderPath []*tile) {
	pathFields := make(map[string]bool)
	for _, t := range renderPath {
		pathFields[fmt.Sprintf("%d:%d", t.x, t.y)] = true
	}

	for y := 0; y <= g.MaxY(); y++ {
		for x := 0; x <= g.MaxX(); x++ {
			if pathFields[fmt.Sprintf("%d:%d", x, y)] {
				fmt.Fprintf(w, "·")
				continue
			}

			fmt.Fprintf(w, "%s", string(byte(g.ElevationAt(x, y))))
		}
		fmt.Fprintln(w)
	}
}

func (g grid) ShortestPath(from *tile) ([]*tile, float64, bool) {
	path, dist, found := astar.Path(from, g.TileAt(g.target[0], g.target[1]))

	var tilePath []*tile
	for _, pe := range path {
		tilePath = append(tilePath, pe.(*tile))
	}

	return tilePath, dist, found
}

func (g grid) ShortestPathLen(from *tile) (int, bool) {
	p, _, found := g.ShortestPath(from)
	return len(p) - 1, found
}

func (t tile) PathNeighbors() (out []astar.Pather) {
	if t.x > 0 && t.grid.TileAt(t.x-1, t.y).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x-1, t.y))
	}
	if t.x < t.grid.MaxX() && t.grid.TileAt(t.x+1, t.y).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x+1, t.y))
	}
	if t.y > 0 && t.grid.TileAt(t.x, t.y-1).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x, t.y-1))
	}
	if t.y < t.grid.MaxY() && t.grid.TileAt(t.x, t.y+1).elevation <= t.elevation+1 {
		out = append(out, t.grid.TileAt(t.x, t.y+1))
	}

	return out
}

func (t tile) PathNeighborCost(to astar.Pather) float64 {
	target := to.(*tile)
	return float64(target.elevation) - float64(t.elevation) + 1
}

func (t tile) PathEstimatedCost(to astar.Pather) float64 {
	target := to.(*tile)
	return float64(t.manhattenDist(target))
}

func (t tile) manhattenDist(target *tile) int {
	return int(math.Abs(float64(target.x)-float64(t.x)) + math.Abs(float64(target.y)-float64(t.y)))
}

func main() {
	g := readGrid(os.Stdin)

	solution1, _, _ := g.ShortestPath(g.TileAt(g.start[0], g.start[1]))
	fmt.Printf("Solution 1: %d\n", len(solution1)-1)

	g.Print(os.Stdout, solution1)

	solution2Start, solution2 := g.FindShortestACoordPathLen()
	fmt.Printf("Solution 2: %d\n", solution2)

	solution2Path, _, _ := g.ShortestPath(solution2Start)
	g.Print(os.Stdout, solution2Path)
}
