package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstMarkerEndPos(t *testing.T) {
	for cs, pos := range map[string]int{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb":    7,
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      5,
		"nppdvjthqldpwncqszvftbrmjlhg":      6,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 10,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  11,
	} {
		assert.Equal(t, pos, commsStream(cs).GetFirstMarkerEndPos(4), "marker pos of %q", cs)
	}

	for cs, pos := range map[string]int{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb":    19,
		"bvwbjplbgvbhsrlpgdmjqwftvncz":      23,
		"nppdvjthqldpwncqszvftbrmjlhg":      23,
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg": 29,
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw":  26,
	} {
		assert.Equal(t, pos, commsStream(cs).GetFirstMarkerEndPos(14), "marker pos of %q", cs)
	}
}
