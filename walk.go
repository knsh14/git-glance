package main

import (
	"fmt"
	"github.com/libgit2/git2go"
	"sort"
)

func ShowDiff(commit *git.Commit) bool {
	ClearScreen()
	fmt.Println(commit.Message())

	diff := GetDiff(commit, "")
	deltas, _ := diff.NumDeltas()
	choices := map[string]string{
		"n": "next",
		"q": "quit",
	}
	keys := []string{"n", "q"}
	for i := 0; i < deltas; i++ {
		delta, _ := diff.GetDelta(i)
		nf := delta.NewFile
		choices[fmt.Sprint(i)] = nf.Path
		keys = append(keys, fmt.Sprint(i))
	}
	diff.Free()

	for {
		command := GetCommand(choices, keys)

		switch command {
		case "n":
			return true
		case "q":
			return false
		default:
			if v, ok := choices[command]; ok {
				diff := GetDiff(commit, v)
				ds := NewDiffStore()
				diff.ForEach(ds.DiffFileCallBackFunc, git.DiffDetailLines)
				diff.Free()
				ClearScreen()
				_ = ds.PassToLess()
			} else {
				return false
			}
		}
		ClearScreen()
	}
}

func GetDiff(c *git.Commit, path string) *git.Diff {
	t, err := c.Tree()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	ddo, _ := git.DefaultDiffOptions()
	if path != "" {
		ddo.Pathspec = []string{path}
	}
	parent := c.Parent(0)
	ptree, _ := parent.Tree()
	r := c.Owner()
	diff, err := r.DiffTreeToTree(ptree, t, &ddo)
	if err != nil {
		return nil
	}
	return diff
}

func GetCommand(choices map[string]string, keys []string) string {
	fmt.Println("You can choose in below")

	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("\t%s: %s\n", k, choices[k])
	}
	for {
		fmt.Print("Your Choice Is: ")
		var command string
		fmt.Scanln(&command)
		if command == "" {
			fmt.Println("input a command! try again")
			continue
		}
		for k, v := range choices {
			if k == command || v == command {
				return k
			}
		}
	}
}
