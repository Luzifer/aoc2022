package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisplayRender(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	c := newCPU(parseInstructionSet(f))

	buf := new(bytes.Buffer)
	display{c}.Render(buf)

	assert.Equal(t, strings.Join([]string{
		"##..##..##..##..##..##..##..##..##..##..",
		"###...###...###...###...###...###...###.",
		"####....####....####....####....####....",
		"#####.....#####.....#####.....#####.....",
		"######......######......######......####",
		"#######.......#######.......#######.....",
	}, "\n"), strings.TrimSpace(buf.String()))
}

func TestLongSignalStrength(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test data")
	t.Cleanup(func() { f.Close() })

	c := newCPU(parseInstructionSet(f))

	c.SetTick(20)
	assert.Equal(t, 21, c.GetRegisterX(), "During the 20th cycle, register X has the value 21")
	assert.Equal(t, 420, c.GetSignalStrength(), "so the signal strength is 20 * 21 = 420.")

	c.SetTick(60)
	assert.Equal(t, 19, c.GetRegisterX(), "During the 60th cycle, register X has the value 19")
	assert.Equal(t, 1140, c.GetSignalStrength(), "so the signal strength is 60 * 19 = 1140")

	c.SetTick(100)
	assert.Equal(t, 18, c.GetRegisterX(), "During the 100th cycle, register X has the value 18")
	assert.Equal(t, 1800, c.GetSignalStrength(), "so the signal strength is 100 * 18 = 1800")

	c.SetTick(140)
	assert.Equal(t, 21, c.GetRegisterX(), "During the 140th cycle, register X has the value 21")
	assert.Equal(t, 2940, c.GetSignalStrength(), "so the signal strength is 140 * 21 = 2940")

	c.SetTick(180)
	assert.Equal(t, 16, c.GetRegisterX(), "During the 180th cycle, register X has the value 16")
	assert.Equal(t, 2880, c.GetSignalStrength(), "so the signal strength is 180 * 16 = 2880")

	c.SetTick(220)
	assert.Equal(t, 18, c.GetRegisterX(), "During the 220th cycle, register X has the value 18")
	assert.Equal(t, 3960, c.GetSignalStrength(), "so the signal strength is 220 * 18 = 3960")
}

func TestShort(t *testing.T) {
	c := newCPU([]instruction{
		noop{},
		addx{3},
		addx{-5},
	})

	assert.Equal(t, 1, c.GetRegisterX(), "The CPU has a single register, X, which starts with the value 1.")

	c.SetTick(1)
	assert.Equal(t, 1, c.GetRegisterX(), "After the first cycle, the noop instruction finishes execution, doing nothing.")

	c.SetTick(2)
	assert.Equal(t, 1, c.GetRegisterX(), "During the second cycle, X is still 1.")

	c.SetTick(3)
	assert.Equal(t, 1, c.GetRegisterX(), "During the third cycle, X is still 1")

	c.SetTick(4)
	assert.Equal(t, 4, c.GetRegisterX(), "During the fourth cycle, X is still 4.")

	c.SetTick(6)
	assert.Equal(t, -1, c.GetRegisterX(), "After the fifth cycle, the addx -5 instruction finishes execution, setting X to -1.")
}
