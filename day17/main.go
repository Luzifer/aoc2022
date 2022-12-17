package main

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

type (
	chamber struct {
		jetDirs []jetDirection
		jetPtr  int
		width   int
		map2d   map[string]bool
		maxY    int

		curShape    *shape
		curShapePos coord
		shapePtr    int
	}

	coord struct{ x, y int }

	jetDirection byte

	shape [][]byte
)

const (
	jetDirectionLeft  jetDirection = '<'
	jetDirectionRight jetDirection = '>'
)

var (
	shapeBlock = &shape{{'#', '#'}, {'#', '#'}}
	shapeI     = &shape{{'#'}, {'#'}, {'#'}, {'#'}}
	shapeInvL  = &shape{{'#', '#', '#'}, {' ', ' ', '#'}, {' ', ' ', '#'}}
	shapeMinus = &shape{{'#', '#', '#', '#'}}
	shapePlus  = &shape{{' ', '#', ' '}, {'#', '#', '#'}, {' ', '#', ' '}}
	shapeOrder = []*shape{shapeMinus, shapePlus, shapeInvL, shapeI, shapeBlock}
)

func newChamber() *chamber {
	return &chamber{
		map2d: make(map[string]bool),
		maxY:  -1,
	}
}

func (c chamber) Draw(w io.Writer) {
	for y := c.maxY + 3; y >= 0; y-- {
		fmt.Fprintf(w, "|")
		for x := 0; x < c.width; x++ {
			if c.curShape != nil &&
				x >= c.curShapePos.x && x < c.curShapePos.x+len((*c.curShape)[0]) &&
				y >= c.curShapePos.y && y < c.curShapePos.y+len(*c.curShape) &&
				(*c.curShape)[y-c.curShapePos.y][x-c.curShapePos.x] == '#' {
				fmt.Fprintf(w, "@")
				continue
			}

			if c.map2d[coord{x, y}.String()] {
				fmt.Fprintf(w, "#")
			} else {
				fmt.Fprintf(w, "\u00b7")
			}
		}
		fmt.Fprintln(w, "|")
	}
}

func (c *chamber) Tick() (hasSettled bool) {
	if c.curShape == nil {
		c.curShape = shapeOrder[c.shapePtr]
		c.curShapePos = coord{2, c.maxY + 4}

		// Select next shape
		c.shapePtr++
		if c.shapePtr == len(shapeOrder) {
			c.shapePtr = 0
		}
	}

	// First apply jet-stream movement if possible
	if !c.curShape.DoesCollide(c, c.jetDirs[c.jetPtr].MoveFn(c.curShapePos)) {
		c.curShapePos = c.jetDirs[c.jetPtr].MoveFn(c.curShapePos)
	}

	c.jetPtr++
	if c.jetPtr == len(c.jetDirs) {
		c.jetPtr = 0
	}

	// Afterwards apply downwards movement
	if !c.curShape.DoesCollide(c, coord{c.curShapePos.x, c.curShapePos.y - 1}) {
		c.curShapePos.y--
		return false
	}

	c.settleShape()
	c.curShape = nil
	return true
}

func (c chamber) genLine() []byte {
	line := make([]byte, c.width)
	for i := range line {
		line[i] = ' '
	}
	return line
}

func (c *chamber) settleShape() {
	for y := 0; y < len(*c.curShape); y++ {
		for x := 0; x < len((*c.curShape)[y]); x++ {

			drawX := x + c.curShapePos.x
			drawY := y + c.curShapePos.y

			if (*c.curShape)[y][x] == '#' {
				c.map2d[coord{drawX, drawY}.String()] = true
			}

			if drawY > c.maxY {
				c.maxY = drawY
			}
		}
	}
}

func (c coord) String() string { return fmt.Sprintf("%d:%d", c.x, c.y) }

func jetDirsFromBytes(in []byte) (out []jetDirection) {
	out = make([]jetDirection, len(in))
	for i := range in {
		out[i] = jetDirection(in[i])
	}
	return out
}

func (j jetDirection) MoveFn(c coord) coord {
	switch j {
	case jetDirectionLeft:
		return coord{c.x - 1, c.y}

	case jetDirectionRight:
		return coord{c.x + 1, c.y}

	default:
		panic("unknown movement")
	}
}

func (s shape) DoesCollide(c *chamber, pos coord) bool {
	if pos.x < 0 || pos.y < 0 || pos.x+len(s[0]) > c.width {
		// Out of bounds
		return true
	}

	for y := 0; y < len(s); y++ {
		for x := 0; x < len(s[y]); x++ {
			switch s[y][x] {
			case ' ':
				// Cannot collide as it is air

			case '#':
				if c.map2d[coord{x + pos.x, y + pos.y}.String()] {
					// Rock would hit rock: Collides!
					return true
				}

			}
		}
	}

	return false // Nothin collided: Yay
}

func main() {
	jets, _ := io.ReadAll(os.Stdin)
	c := newChamber()
	c.jetDirs = jetDirsFromBytes(bytes.TrimSpace(jets))
	c.width = 7

	fmt.Printf("Solution 1: %d\n", heightAfter(c, 2022))

	c2 := newChamber()
	c2.width = 7
	c2.jetDirs = c.jetDirs

	fmt.Printf("Solution 2: %d\n", heightAfter(c, 1000000000000))
}

func heightAfter(c *chamber, n int64) int {
	p := pb.New64(int64(n))
	p.Start()
	defer p.Finish()

	var settled int64
	for settled < n {
		if c.Tick() {
			settled++
			p.Set64(settled)
		}
	}

	return c.maxY + 1
}
