package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

func gitinfo(w io.Writer) {
	if _, err := os.Stat("./.git"); errors.Is(err, fs.ErrNotExist) {
		return
	}

	head, err := os.Open("./.git/HEAD")
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
