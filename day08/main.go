package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

type (
	tree struct {
		Height int
		X, Y   int
	}
	grid [][]tree
)

func readGrid(r io.Reader) grid {
	var (
		g       grid
		scanner = bufio.NewScanner(r)
	)

	for scanner.Scan() {
		var row []tree
		for i, c := range scanner.Bytes() {
			tree := tree{X: i, Y: len(g)}
			tree.Height, _ = strconv.Atoi(string(c))
			row = append(row, tree)
		}
		g = append(g, row)
	}

	return g
}

func (g grid) GetTree(x, y int) tree { return g[y][x] }

func (g grid) IterateTrees(f func(tree)) {
	for y := 0; y <= g.MaxY(); y++ {
		for x := 0; x <= g.MaxX(); x++ {
			f(g.GetTree(x, y))
		}
	}
}

func (g grid) MaxX() int {
	if len(g) == 0 {
		return 0
	}
	return len(g[0]) - 1
}

func (g grid) MaxY() int { return len(g) - 1 }

func (t tree) GetScenicScore(g grid) int {
	score := 1

	for _, modFn := range []func(x, y int) (int, int){
		func(x, y int) (int, int) { return x, y - 1 },
		func(x, y int) (int, int) { return x, y + 1 },
		func(x, y int) (int, int) { return x - 1, y },
		func(x, y int) (int, int) { return x + 1, y },
	} {
		var (
			coordX, coordY = t.X, t.Y
			seenTrees      int
		)
		for {
			coordX, coordY = modFn(coordX, coordY)
			if coordX < 0 || coordX > g.MaxX() || coordY < 0 || coordY > g.MaxY() {
				// We ran outside the field, stop looking
				break
			}

			seenTrees++
			if g.GetTree(coordX, coordY).Height >= t.Height {
				// That's the last one we can see
				break
			}
		}

		score *= seenTrees
	}

	return score
}

func (t tree) IsOnEdge(g grid) bool {
	return t.X == 0 || t.X == g.MaxX() ||
		t.Y == 0 || t.Y == g.MaxY()
}

func (t tree) IsOutsideVisible(g grid) bool {
	if t.IsOnEdge(g) {
		return true
	}

	for _, coordSet := range [][2][]int{
		{{t.X}, t.mkSeq(0, t.Y-1)},
		{{t.X}, t.mkSeq(t.Y+1, g.MaxY())},
		{t.mkSeq(0, t.X-1), {t.Y}},
		{t.mkSeq(t.X+1, g.MaxX()), {t.Y}},
	} {
		setVisible := true
		for _, x := range coordSet[0] {
			for _, y := range coordSet[1] {
				if g.GetTree(x, y).Height >= t.Height {
					setVisible = false
				}
			}
		}

		if setVisible {
			return true
		}
	}

	return false
}

func (t tree) mkSeq(min, max int) []int {
	var out []int
	for i := min; i <= max; i++ {
		out = append(out, i)
	}

	return out
}

func main() {
	g := readGrid(os.Stdin)

	var solution1, solution2 int
	g.IterateTrees(func(t tree) {
		if t.IsOutsideVisible(g) {
			solution1++
		}

		if s := t.GetScenicScore(g); s > solution2 {
			solution2 = s
		}
	})

	fmt.Printf("Solution 1: %d\n", solution1)
	fmt.Printf("Solution 2: %d\n", solution2)
}
