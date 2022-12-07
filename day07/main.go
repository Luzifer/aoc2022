package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type (
	fsEntry struct {
		IsDir      bool
		Name       string
		Parent     *fsEntry
		Size       uint64
		SubEntries []*fsEntry
	}
)

func readFSStructure(r io.Reader) *fsEntry {
	rootFSEntry := &fsEntry{IsDir: true}

	var (
		curDir  string
		scanner = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), " ")
		if parts[0] == "$" {
			// Command
			switch parts[1] {
			case "cd":
				switch parts[2] {
				case "/":
					curDir = ""

				case "..":
					cp := strings.Split(curDir, "/")
					curDir = strings.Join(cp[:len(cp)-1], "/")

				default:
					curDir = strings.Join([]string{curDir, parts[2]}, "/")
				}

			case "ls":
				// Whatever

			default:
				panic("unknown command")
			}

			continue
		}

		var newEntry *fsEntry
		// Dir entry
		if parts[0] == "dir" {
			newEntry = &fsEntry{IsDir: true, Name: parts[1]}
		} else {
			// Should be file
			newEntry = &fsEntry{IsDir: false, Name: parts[1]}
			newEntry.Size, _ = strconv.ParseUint(parts[0], 10, 64)
		}

		rootFSEntry.GetByPath(curDir).Add(newEntry)
	}

	return rootFSEntry
}

func (f *fsEntry) Add(e *fsEntry) {
	e.Parent = f
	f.SubEntries = append(f.SubEntries, e)
}

func (f *fsEntry) CalcSize() uint64 {
	if !f.IsDir {
		return f.Size
	}

	var size uint64
	for _, e := range f.SubEntries {
		size += e.CalcSize()
	}

	return size
}

func (f *fsEntry) GetByPath(p string) *fsEntry {
	parts := strings.Split(p, "/")
	switch len(parts) {
	case 0:
		// WTF?
		panic("invalid path")

	case 1:
		if parts[0] == f.Name {
			return f
		}
		return nil

	default:
		if parts[0] != f.Name {
			return nil
		}
		for _, e := range f.SubEntries {
			if r := e.GetByPath(strings.Join(parts[1:], "/")); r != nil {
				return r
			}
		}
	}

	return nil
}

func (f *fsEntry) GetDirs() (out []*fsEntry) {
	if !f.IsDir {
		return nil
	}

	out = append(out, f)

	for _, e := range f.SubEntries {
		if !e.IsDir {
			continue
		}

		out = append(out, e.GetDirs()...)
	}

	return out
}

func (f *fsEntry) HasSubdirs() bool {
	for _, e := range f.SubEntries {
		if e.IsDir {
			return true
		}
	}

	return false
}

func (f *fsEntry) Path() string {
	if f.Parent == nil {
		return f.Name
	}

	return strings.Join([]string{f.Parent.Path(), f.Name}, "/")
}

func main() {
	fs := readFSStructure(os.Stdin)

	// Solution 1
	var (
		dirs                 = fs.GetDirs()
		solution1            uint64
		solution2            uint64 = 70000000
		solution2Requirement        = 30000000 - (70000000 - fs.GetByPath("").CalcSize())
	)
	for _, d := range dirs {
		size := d.CalcSize()

		if size <= 100000 {
			solution1 += size
			continue
		}

		if size < solution2 && size >= solution2Requirement {
			solution2 = size
		}
	}

	fmt.Printf("Solution 1: %d\n", solution1)
	fmt.Printf("Solution 2: %d\n", solution2)
}
