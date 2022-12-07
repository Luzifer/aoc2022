package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadFSStructure(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test file")

	fs := readFSStructure(f)
	require.NotNil(t, fs, "expect to have root node")

	assert.Equal(t, uint64(14848514), fs.GetByPath("/b.txt").Size, "checking size of b.txt")
	assert.Equal(t, uint64(62596), fs.GetByPath("/a/h.lst").Size, "checking size of h.lst")
	assert.Equal(t, uint64(5626152), fs.GetByPath("/d/d.ext").Size, "checking size of d.ext")
}

func TestCalcSize(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test file")

	fs := readFSStructure(f)
	require.NotNil(t, fs, "expect to have root node")

	assert.Equal(t, uint64(584), fs.GetByPath("/a/e").CalcSize(), "checking size of /a/e")
	assert.Equal(t, uint64(94853), fs.GetByPath("/a").CalcSize(), "checking size of /a")
	assert.Equal(t, uint64(24933642), fs.GetByPath("/d").CalcSize(), "checking size of /d")
	assert.Equal(t, uint64(48381165), fs.GetByPath("").CalcSize(), "checking size of /")
}

func TestGetDirs(t *testing.T) {
	f, err := os.Open("test.txt")
	require.NoError(t, err, "opening test file")

	fs := readFSStructure(f)
	require.NotNil(t, fs, "expect to have root node")

	fd := fs.GetDirs()
	assert.Len(t, fd, 4, "collecting dirs")
}
