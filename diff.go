package main

import (
	"bytes"
	"fmt"
	"github.com/libgit2/git2go"
	"os"
	"os/exec"
)

func NewDiffStore() *DiffStore {
	d := &DiffStore{diffs: make([]string, 1)}
	return d
}

type DiffStore struct {
	diffs []string
}

func (d *DiffStore) DiffLineCallBackFunc(dl git.DiffLine) error {
	switch dl.Origin {
	case git.DiffLineAddition:
		d.Store(fmt.Sprintf("\x1b[32m+ %s\x1b[0m", dl.Content))
	case git.DiffLineDeletion:
		d.Store(fmt.Sprintf("\x1b[31m- %s\x1b[0m", dl.Content))
	case git.DiffLineContext:
		d.Store(fmt.Sprintf(" %s", dl.Content))
	}
	return nil
}

func (d *DiffStore) DiffHunkCallBackFunc(dh git.DiffHunk) (git.DiffForEachLineCallback, error) {
	d.Store(fmt.Sprintln(dh.Header))
	return d.DiffLineCallBackFunc, nil
}

func (d *DiffStore) DiffFileCallBackFunc(dd git.DiffDelta, n float64) (git.DiffForEachHunkCallback, error) {
	d.Store(fmt.Sprintf("%s -> %s\n", dd.OldFile.Path, dd.NewFile.Path))
	return d.DiffHunkCallBackFunc, nil
}

func (d *DiffStore) PassToLess() error {
	b := StringsToBytes(d.diffs)
	cmd := exec.Command("less", "-r")
	cmd.Stdin = bytes.NewReader(b)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (d *DiffStore) Store(s string) {
	d.diffs = append(d.diffs, s)
}

func StringsToBytes(diffs []string) []byte {
	b := make([]byte, 0)
	for _, s := range diffs {
		t := []byte(s)
		b = append(b, t...)
	}
	return b
}
