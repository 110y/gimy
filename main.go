package main

import (
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "not enough arguments")
		os.Exit(1)
	}

	path := os.Args[1]

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse file: %s\n", err)
		os.Exit(1)
	}

	for _, im := range f.Imports {
		if im.Name != nil {
			im.Name.NamePos = token.NoPos
		}
		im.Path.ValuePos = token.NoPos
		im.EndPos = token.NoPos
	}

	tmpfile, err := ioutil.TempFile("", "gimy-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create tmp file: %s\n", err)
		os.Exit(1)
	}
	defer tmpfile.Close()

	if err := format.Node(tmpfile, fset, f); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write file: %s\n", err)
		os.Exit(1)
	}

	if err := os.Rename(tmpfile.Name(), path); err != nil {
		fmt.Fprintf(os.Stderr, "failed to rename file: %s\n", err)
		os.Exit(1)
	}
}
