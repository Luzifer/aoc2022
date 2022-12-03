package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

func prioForItem(item byte) uint {
	if 'a' <= item && item <= 'z' {
		return uint(item) - uint('a') + 1
	}
	return uint(item) - uint('A') + 27
}

func main() {
	var (
		solution1, solution2 uint
		groupMembers         [][]byte
	)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}

		if len(scanner.Text())%2 != 0 {
			panic("Invalid input")
		}

		compartment1 := scanner.Bytes()[:len(scanner.Bytes())/2]
		compartment2 := scanner.Bytes()[len(scanner.Bytes())/2:]

		var accountedFor []byte
		for _, item := range compartment1 {
			if bytes.Contains(compartment2, []byte{item}) && !bytes.Contains(accountedFor, []byte{item}) {
				solution1 += prioForItem(item)
				accountedFor = append(accountedFor, item)
			}
		}

		groupMembers = append(groupMembers, scanner.Bytes())
		if len(groupMembers) == 3 {
			var accountedFor []byte
			for _, item := range groupMembers[0] {
				if bytes.Contains(accountedFor, []byte{item}) {
					// Duplicate, don't count again
					continue
				}

				if bytes.Contains(groupMembers[1], []byte{item}) && bytes.Contains(groupMembers[2], []byte{item}) {
					// All three carry this item, this must be the identifier
					solution2 += prioForItem(item)
				}

				accountedFor = append(accountedFor, item)
			}

			groupMembers = nil
		}
	}

	fmt.Printf("Solution 1: %d\n", solution1)
	fmt.Printf("Solution 2: %d\n", solution2)
}
