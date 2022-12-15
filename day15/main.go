package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type (
	coord struct{ x, y int64 }

	sensor struct {
		pos         coord
		beaconFound coord
	}

	sensorArray []*sensor
)

var inputRegex = regexp.MustCompile(`Sensor at x=([0-9-]+), y=([0-9-]+): closest beacon is at x=([0-9-]+), y=([0-9-]+)`)

func readSensorArray(r io.Reader) (a sensorArray) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		matches := inputRegex.FindStringSubmatch(scanner.Text())
		s := &sensor{}
		s.pos.x, _ = strconv.ParseInt(matches[1], 10, 64)
		s.pos.y, _ = strconv.ParseInt(matches[2], 10, 64)
		s.beaconFound.x, _ = strconv.ParseInt(matches[3], 10, 64)
		s.beaconFound.y, _ = strconv.ParseInt(matches[4], 10, 64)
		a = append(a, s)
	}

	return a
}

func (s sensor) BeaconDist() int64 {
	return s.manhatten(s.beaconFound)
}

// CoveredFieldsAtY returns only X part of the coordinate
func (s sensor) CoveredFieldsAtY(y int64) []int64 {
	if !s.HasInfoAbout(coord{s.pos.x, y}) {
		// Y coordinate is out of reach, we don't need to scan for
		return nil
	}

	var out []int64
	// Lets scan the line
	x := s.pos.x
	for s.HasInfoAbout(coord{x, y}) {
		out = append(out, x)
		x--
	}

	// Now lets also scan to the right
	x = s.pos.x + 1
	for s.HasInfoAbout(coord{x, y}) {
		out = append(out, x)
		x++
	}

	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (s sensor) HasInfoAbout(c coord) bool {
	return s.manhatten(c) <= s.BeaconDist()
}

func (s sensor) manhatten(c coord) int64 {
	return int64(math.Abs(float64(c.x-s.pos.x)) + math.Abs(float64(c.y-s.pos.y)))
}

func main() {
	a := readSensorArray(os.Stdin)

	fmt.Printf("Solution 1: %d\n", solve1(a, 2000000))
	fmt.Printf("Solution 2: %d\n", solve2(a, 4000000))
}

func solve1(a sensorArray, y int64) int {
	xs := map[int64]bool{}
	for _, s := range a {
		for _, x := range s.CoveredFieldsAtY(y) {
			xs[x] = true
		}
	}

	for _, s := range a {
		if s.beaconFound.y != y {
			continue
		}

		xs[s.beaconFound.x] = false // It's a beacon, do not list it
	}

	var count int
	for _, ok := range xs {
		if ok {
			count++
		}
	}

	return count
}

func solve2(a sensorArray, clim int64) int64 {
	pos := coord{0, 0}
	for {
		var matchingSensor *sensor
		for i := range a {
			if !a[i].HasInfoAbout(pos) {
				continue
			}
			matchingSensor = a[i]
			break
		}

		if matchingSensor == nil {
			return pos.x*4000000 + pos.y
		}

		// We can skip ahead and ignore some fields we don't need to check
		pos.x += matchingSensor.BeaconDist() - matchingSensor.manhatten(pos) + 1

		if pos.x > clim {
			pos.x = 0
			pos.y++
		}

		if pos.y > clim {
			break
		}
	}

	panic("unocovered")
}
