package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type (
	cargo struct {
		stacks [][]byte
	}

	moveMode uint
)

const (
	moveMode9000 moveMode = iota
	moveMode9001
)

func readCargo(scanner *bufio.Scanner) (*cargo, error) {
	var c cargo

	for scanner.Scan() {
		if scanner.Text() == "" {
			// Cargo table is divided from move set by one empty line
			break
		}

		line := append(scanner.Bytes(), ' ')
		for i := 0; i < len(line); i += 4 {
			if c.stacks == nil {
				c.stacks = make([][]byte, len(line)/4)
			}

			if '1' <= line[1] && line[1] <= '9' {
				// Number line, we don't care
				continue
			}

			// As we're reading the stack from top to bottom but expect the
			// top to be the last item in the slice we need to prepend the
			// newly read crate identifier
			if line[i+1] != ' ' {
				c.stacks[i/4] = append([]byte{line[i+1]}, c.stacks[i/4]...)
			}
		}
	}

	return &c, nil
}

func (c cargo) Clone() *cargo {
	n := &cargo{
		stacks: make([][]byte, len(c.stacks)),
	}

	for s := 0; s < len(c.stacks); s++ {
		for i := 0; i < len(c.stacks[s]); i++ {
			n.stacks[s] = append(n.stacks[s], c.stacks[s][i])
		}
	}

	return n
}

func (c *cargo) GetTopCrate(stack int) byte {
	// Puzzle works 1-based, we are working 0-based
	stack -= 1

	return c.stacks[stack][len(c.stacks[stack])-1]
}

func (c *cargo) Len() int {
	return len(c.stacks)
}

func (c *cargo) Move(from, to, count int, mode moveMode) {
	// Puzzle works 1-based, we are working 0-based
	from, to = from-1, to-1

	switch mode {
	case moveMode9000:
		for i := 0; i < count; i++ {
			// Add last element (top) from [from] to [to]
			c.stacks[to] = append(c.stacks[to], c.stacks[from][len(c.stacks[from])-1])
			// Remove last element from [from]
			c.stacks[from] = c.stacks[from][:len(c.stacks[from])-1]
		}

	case moveMode9001:
		c.stacks[to] = append(c.stacks[to], c.stacks[from][len(c.stacks[from])-count:]...)
		c.stacks[from] = c.stacks[from][:len(c.stacks[from])-count]
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	c, err := readCargo(scanner)
	if err != nil {
		panic(err)
	}

	c2 := c.Clone()

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		count, _ := strconv.Atoi(parts[1])
		from, _ := strconv.Atoi(parts[3])
		to, _ := strconv.Atoi(parts[5])

		c.Move(from, to, count, moveMode9000)
		c2.Move(from, to, count, moveMode9001)
	}

	var solution1, solution2 []byte
	for i := 1; i <= c.Len(); i++ {
		solution1 = append(solution1, c.GetTopCrate(i))
		solution2 = append(solution2, c2.GetTopCrate(i))
	}

	fmt.Printf("Solution 1: %s\n", solution1)
	fmt.Printf("Solution 2: %s\n", solution2)
}
