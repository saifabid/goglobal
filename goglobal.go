package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	flag.Usage = usage
	flag.Parse()
	os.Exit(run(os.Args))
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: goglobal [path ...]\n")
	flag.PrintDefaults()
	os.Exit(0)
}

func run(args []string) int {
	if len(args) == 1 {
		fmt.Println("no path provided")
		return 0
	}

	p := args[1]
	if p == "" {
		fmt.Println("invalid path")
		return 0
	}

	switch dir, err := os.Stat(p); {
	case err != nil:
		fmt.Println(err)
		return 1
	case dir.IsDir():
		return walk(p)
	default:
		fmt.Println("path unreadable")
		return 0
	}
}

func walk(path string) int {
	err := filepath.Walk(path, work)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func work(path string, info os.FileInfo, err error) error {
	if !isGoFile(info) {
		return nil
	}

	fset := token.NewFileSet()
	fp, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}

	for _, s := range fp.Scope.Objects {
		if s.Kind.String() == "var" {
			fmt.Printf("%s %d:\t%s\n", path, fset.File(s.Pos()).Line(s.Pos()), s.Name)
		}
	}

	return nil
}

func isGoFile(f os.FileInfo) bool {
	return !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && strings.HasSuffix(f.Name(), ".go")
}
