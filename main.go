package main

import (
	"fmt"
	"github.com/libgit2/git2go"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Printf("Usage: %s <ref>\n", os.Args[0])
}

func main() {
	args := os.Args
	if len(args) < 2 || len(args) > 2 {
		usage()
		return
	}
	ref := args[1]

	p, _ := filepath.Abs(".")
	repo, err := git.OpenRepository(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	walk, err := repo.Walk()
	if err != nil {
		fmt.Println(err)
		return
	}
	commitWay := fmt.Sprintf("%s..HEAD", ref)
	err = walk.PushRange(commitWay)
	if err != nil {
		fmt.Println(err)
		usage()
		return
	}
	walk.Sorting(git.SortReverse)

	walk.Iterate(ShowDiff)
	ClearScreen()
}

func ClearScreen() {
	// clear display before quit
	fmt.Print("\033[2J")
	fmt.Printf("\033[%d;%dH", 0, 0)
}
