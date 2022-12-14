package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type (
	coord struct{ x, y int }

	fieldType uint

	grid struct {
		fields                 map[string]fieldType
		maxX, maxY, minX, minY int
		movingSand             *coord
	}
)

const (
	fieldTypeAir fieldType = iota
	fieldTypeRock
	fieldTypeSand
)

func (f fieldType) String() string {
	return map[fieldType]string{
		fieldTypeAir:  "air",
		fieldTypeRock: "rock",
		fieldTypeSand: "sand",
	}[f]
}

var (
	errSandFellToVoid   = errors.New("sand vanished to the void")
	errSandInputClogged = errors.New("sand clogged the input")
	errSandSettled      = errors.New("sand settled and cannot move anymore")
)

func newGrid() *grid {
	return &grid{
		fields: make(map[string]fieldType),
		minX:   math.MaxInt,
		minY:   0, // Forcing to zero as 500,0 is the source of the sand
	}
}

func (g grid) MaterialAt(x, y int) fieldType { return g.fields[g.ctos(x, y)] }
func (g grid) MaterialAtWithSimulatedFloor(floorY, x, y int) fieldType {
	if y == floorY {
		return fieldTypeRock
	}
	return g.MaterialAt(x, y)
}

func (g *grid) MoveSandGrain(simulateFloor bool) error {
	if g.movingSand == nil {
		if g.MaterialAt(500, 0) == fieldTypeSand {
			// We cannot spawn sand if there is already sand
			return errSandInputClogged
		}

		g.movingSand = &coord{500, 0} // Spawn sand!
	}

	if g.movingSand.y > g.maxY+10 {
		// There is nothing more according to our notes so baaaaaji sand!
		return errSandFellToVoid
	}

	fieldCheck := g.MaterialAt
	if simulateFloor {
		fieldCheck = func(x, y int) fieldType { return g.MaterialAtWithSimulatedFloor(g.maxY+2, x, y) }
	}

	switch {
	case fieldCheck(g.movingSand.x, g.movingSand.y+1) == fieldTypeAir:
		g.movingSand.y++
		return nil

	case fieldCheck(g.movingSand.x-1, g.movingSand.y+1) == fieldTypeAir:
		g.movingSand.y++
		g.movingSand.x--
		return nil

	case fieldCheck(g.movingSand.x+1, g.movingSand.y+1) == fieldTypeAir:
		g.movingSand.y++
		g.movingSand.x++
		return nil

	default:
		// Sand cannot move and is settled now
		g.fields[g.ctos(g.movingSand.x, g.movingSand.y)] = fieldTypeSand
		g.movingSand = nil
		return errSandSettled
	}
}

func (g *grid) ReadInput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		g.ReadPath(scanner.Text())
	}
}

func (g *grid) ReadPath(line string) {
	var path [][2]int
	coordStrs := strings.Split(line, " -> ")
	for _, coordStr := range coordStrs {
		coordParts := strings.Split(coordStr, ",")
		x, _ := strconv.Atoi(coordParts[0])
		y, _ := strconv.Atoi(coordParts[1])
		path = append(path, [2]int{x, y})
	}

	for _, coord := range g.extrapolatePath(path) {
		g.fields[g.ctos(coord[0], coord[1])] = fieldTypeRock
		if coord[0] > g.maxX {
			g.maxX = coord[0]
		}
		if coord[0] < g.minX {
			g.minX = coord[0]
		}
		if coord[1] > g.maxY {
			g.maxY = coord[1]
		}
		if coord[1] < g.minY {
			g.minY = coord[1]
		}
	}
}

func (g grid) Render(w io.Writer) {
	fmt.Fprintf(w, "%d,%d : %d,%d\n", g.minX, g.minY, g.maxX, g.maxY)
	for y := g.minY; y <= g.maxY; y++ {
		for x := g.minX; x <= g.maxX; x++ {
			if g.movingSand != nil && x == g.movingSand.x && y == g.movingSand.y {
				fmt.Fprintf(w, "o")
			}

			switch g.MaterialAt(x, y) {
			case fieldTypeAir:
				fmt.Fprintf(w, "\u00b7")
			case fieldTypeRock:
				fmt.Fprintf(w, "#")
			case fieldTypeSand:
				fmt.Fprintf(w, "o")
			}
		}
		fmt.Fprintln(w)
	}
}

func (g *grid) SettleSandGrains(n int, simulateFloor bool) (int, error) {
	var sandsSettled int

	for i := 0; i < n; i++ {
		for {
			err := g.MoveSandGrain(simulateFloor)
			switch {
			case err == nil:
				// Sand has moved and needs to move again
				continue

			case errors.Is(err, errSandSettled):
				// That grain of sand will never move again
				sandsSettled++
				break

			case errors.Is(err, errSandFellToVoid) || errors.Is(err, errSandInputClogged):
				return sandsSettled, err

			default:
				panic("WTF?")
			}

			break
		}
	}

	return sandsSettled, nil
}

func (grid) ctos(x, y int) string { return fmt.Sprintf("%d:%d", x, y) }

func (grid) extrapolatePath(path [][2]int) [][2]int {
	var (
		outPath [][2]int
		lastC   = path[0]
	)
	for i := 1; i < len(path); i++ {
		c := path[i]
		switch {
		case c[0] < lastC[0]:
			for x := c[0]; x <= lastC[0]; x++ {
				outPath = append(outPath, [2]int{x, c[1]})
			}

		case c[0] > lastC[0]:
			for x := lastC[0]; x <= c[0]; x++ {
				outPath = append(outPath, [2]int{x, c[1]})
			}

		case c[1] < lastC[1]:
			for y := c[1]; y <= lastC[1]; y++ {
				outPath = append(outPath, [2]int{c[0], y})
			}

		case c[1] > lastC[1]:
			for y := lastC[1]; y <= c[1]; y++ {
				outPath = append(outPath, [2]int{c[0], y})
			}
		}
		lastC = c
	}

	return outPath
}

func main() {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, os.Stdin); err != nil {
		panic(err)
	}

	// ----

	g := newGrid()
	g.ReadInput(bytes.NewReader(buf.Bytes()))

	var solution1 int
	for {
		n, err := g.SettleSandGrains(500, false)
		solution1 += n
		if errors.Is(err, errSandFellToVoid) {
			break
		}
	}
	fmt.Printf("Solution 1: %d\n", solution1)
	g.Render(os.Stderr)

	// ----

	g = newGrid()
	g.ReadInput(bytes.NewReader(buf.Bytes()))

	var solution2 int
	for {
		n, err := g.SettleSandGrains(500, true)
		solution2 += n
		if errors.Is(err, errSandInputClogged) {
			break
		}
	}
	fmt.Printf("Solution 2: %d\n", solution2)
	g.Render(os.Stderr)
}
