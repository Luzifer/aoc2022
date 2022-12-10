package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	crtWidth  = 40
	crtHeight = 6
)

type (
	cpu struct {
		instructionSet []instruction
		tick, x        int
	}

	display struct {
		c *cpu
	}

	instruction interface {
		ModifyRegisterX(int) int
		TickCount() int
	}

	addx struct{ mod int }
	noop struct{}
)

var ( // Compile-time assertions
	_ instruction = addx{}
	_ instruction = noop{}
)

func parseInstructionSet(r io.Reader) (is []instruction) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		switch parts[0] {
		case "addx":
			v, _ := strconv.Atoi(parts[1])
			is = append(is, addx{v})

		case "noop":
			is = append(is, noop{})
		}
	}

	return is
}

func (a addx) ModifyRegisterX(x int) int { return x + a.mod }
func (addx) TickCount() int              { return 2 }

func (noop) ModifyRegisterX(x int) int { return x }
func (noop) TickCount() int            { return 1 }

// newCPU returns a CPU primed with the given instruction set at tick 0
func newCPU(is []instruction) *cpu { return &cpu{instructionSet: is, tick: 1, x: 1} }

// GetRegisterX returns the current value of the register X
func (c cpu) GetRegisterX() int { return c.x }

// GetSignalStrength returns the signal strength at the current tick
func (c cpu) GetSignalStrength() int { return c.tick * c.x }

// SetTick puts the CPU at the state of the given tick number
func (c *cpu) SetTick(tick int) {
	var currentInst, ticksInCurrentInst int

	c.tick, c.x = 0, 1

	for i := 0; i < tick; i++ {
		if ticksInCurrentInst == c.instructionSet[currentInst].TickCount() {
			// Instruction was finished with this tick
			c.x = c.instructionSet[currentInst].ModifyRegisterX(c.x) // Apply instruction effect

			currentInst++
			ticksInCurrentInst = 0
		}

		c.tick++
		ticksInCurrentInst++ // Tick is finished, we increase by one
	}
}

func (d display) Render(o io.Writer) {
	for px := 0; px < crtHeight*crtWidth; px++ {
		if px%crtWidth == 0 {
			fmt.Fprintln(o)
		}

		// Each tick one pixel is drawn during that tick, during the first
		// tick the first (0) pixel is drawn
		d.c.SetTick(px + 1)

		pxInRow := px % crtWidth
		if x := d.c.GetRegisterX(); x-1 <= pxInRow && pxInRow <= x+1 {
			fmt.Fprintf(o, "#")
		} else {
			fmt.Fprintf(o, ".")
		}
	}
	fmt.Fprintln(o)
}

func main() {
	c := newCPU(parseInstructionSet(os.Stdin))

	var solution1 int
	for _, tick := range []int{20, 60, 100, 140, 180, 220} {
		c.SetTick(tick)
		solution1 += c.GetSignalStrength()
	}

	fmt.Printf("Solution 1: %d\n", solution1)

	fmt.Printf("Solution 2:")
	display{c}.Render(os.Stdout)
}
