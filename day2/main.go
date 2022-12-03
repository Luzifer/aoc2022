package main

import (
	"bufio"
	"fmt"
	"os"
)

type (
	result uint
	shape  uint
)

const (
	resultLoss result = iota
	resultDraw
	resultWin
)

const (
	shapeRock shape = iota + 1
	shapePaper
	shapeScissors
)

func resultFromInputChar(c byte) result {
	switch c {
	case 'X':
		return resultLoss
	case 'Y':
		return resultDraw
	case 'Z':
		return resultWin

	default:
		panic("Invalid input char")
	}
}

func (r result) Score() uint {
	switch r {
	case resultLoss:
		return 0
	case resultDraw:
		return 3
	case resultWin:
		return 6

	default:
		panic("Invalid result")
	}
}

func shapeFromInputChar(c byte) shape {
	switch c {
	case 'A', 'X':
		return shapeRock
	case 'B', 'Y':
		return shapePaper
	case 'C', 'Z':
		return shapeScissors

	default:
		panic("Invalid input char")
	}
}

func (s shape) GetShapeForResult(r result) shape {
	if r == resultDraw {
		// Easy, no need to waste time in switches
		return s
	}

	switch s {
	case shapeRock:
		switch r {
		case resultLoss:
			return shapeScissors
		case resultWin:
			return shapePaper
		}

	case shapeScissors:
		switch r {
		case resultLoss:
			return shapePaper
		case resultWin:
			return shapeRock
		}

	case shapePaper:
		switch r {
		case resultLoss:
			return shapeRock
		case resultWin:
			return shapeScissors
		}

	default:
		panic("Invalid shape")
	}

	panic("Invalid combination")
}

func (s shape) ResultAgainst(i shape) result {
	if s == i {
		// Easy: Both have the same, it's a draw
		return resultDraw
	}

	if s == shapeRock && i == shapeScissors ||
		s == shapeScissors && i == shapePaper ||
		s == shapePaper && i == shapeRock {
		// These are the winning combinations
		return resultWin
	}

	// Any other combination must be a loss
	return resultLoss
}

func (s shape) Score() uint { return uint(s) }

func main() {
	var scorePuzzle1, scorePuzzle2 uint

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		opponent := shapeFromInputChar(scanner.Bytes()[0])
		response := shapeFromInputChar(scanner.Bytes()[2])

		scorePuzzle1 += response.Score()
		scorePuzzle1 += response.ResultAgainst(opponent).Score()

		respForExpectation := opponent.GetShapeForResult(resultFromInputChar(scanner.Bytes()[2]))

		scorePuzzle2 += respForExpectation.Score()
		scorePuzzle2 += respForExpectation.ResultAgainst(opponent).Score()
	}

	fmt.Printf("Solution 1: %d\n", scorePuzzle1)
	fmt.Printf("Solution 2: %d\n", scorePuzzle2)
}
