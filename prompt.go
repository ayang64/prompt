/*

Prompt is a really simple prompt decorator taylor made to meet my preferences.

*/

package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// gitinfo writes the git branch to the provided writer.
func gitinfo(w io.Writer) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}

	hasGit := func(p string) bool {
		// ya, i know.  i was just dinking around here.  will simplify this later.
		return func() error { _, err := os.Stat(path.Join(p, ".git")); return err }() == nil
	}

	var cur string
	for cur = path.Clean(wd); !hasGit(cur) && cur != "/"; cur = path.Dir(cur) {
	}

	head, err := os.Open(path.Join(cur, "./.git/HEAD"))
	if err != nil {
		return
	}
	defer head.Close()

	var ref string
	matches, err := fmt.Fscanf(head, "ref: %s", &ref)
	if err != nil || matches != 1 {
		return
	}

	fmt.Fprintf(w, "%s", strings.TrimPrefix(ref, "refs/heads/"))
}

func main() {
	prompt := strings.Builder{}

	gitinfo(&prompt)

	if prompt.Len() == 0 {
		return
	}

	fmt.Printf("%s - ", &prompt)
}
