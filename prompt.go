/*

Prompt is a very simple prompt decorator taylor made to meet my preferences.

At the moment it only shows the current Git branch.

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
func gitinfo(w io.Writer) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// iterate over paths until we either reach the root or find a .git
	// directory.
	findGit := func() (string, error) {
		hasGit := func(p string) bool {
			// ya, i know.
			return func() error { _, err := os.Stat(path.Join(p, ".git")); return err }() == nil
		}
		for cur := path.Clean(wd); ; cur = path.Dir(cur) {
			switch {
			case hasGit(cur):
				return cur, nil
			case cur == "/":
				return "", fmt.Errorf(".git directory not found")
			}
		}
	}

	dir, err := findGit()
	if err != nil {
		return err
	}

	head, err := os.Open(path.Join(dir, "./.git/HEAD"))
	if err != nil {
		return err
	}
	defer head.Close()

	var ref string
	matches, err := fmt.Fscanf(head, "ref: %s", &ref)
	if err != nil {
		return err
	}
	if matches != 1 {
		return fmt.Errorf("could not parse HEAD file")
	}

	if _, err := fmt.Fprintf(w, "%s", strings.TrimPrefix(ref, "refs/heads/")); err != nil {
		return err
	}

	return nil
}

func main() {
	prompt := strings.Builder{}
	gitinfo(&prompt)

	if prompt.Len() == 0 {
		return
	}

	fmt.Printf("%s - ", &prompt)
}
