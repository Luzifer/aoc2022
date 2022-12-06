package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type commsStream []byte

func (c commsStream) GetFirstMarkerEndPos(markerLen int) int {
	for i := markerLen; i < len(c); i++ {
		var (
			sample = c[i-markerLen : i]
			hasDup bool
		)
		for j := 0; j < len(sample); j++ {
			if bytes.Count(sample, sample[j:j+1]) > 1 {
				hasDup = true
				break
			}
		}

		if !hasDup {
			return i
		}
	}

	panic("no marker found")
}

func main() {
	rawIn, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	input := commsStream(bytes.TrimRight(rawIn, "\n"))

	fmt.Printf("Solution 1: %d\n", input.GetFirstMarkerEndPos(4))
	fmt.Printf("Solution 2: %d\n", input.GetFirstMarkerEndPos(14))
}
