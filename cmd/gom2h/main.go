package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/matsuyoshi30/gom2h"
)

const name = "gom2h"

func main() {
	os.Exit(run(os.Args[1:]))
}

const (
	exitOK = iota
	exitNG
)

var (
	tmpl1 = []byte(`<!DOCTYPE html>
<html>
<head>
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/styles/agate.min.css">
<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/9.18.1/highlight.min.js"></script>
<script>hljs.initHighlightingOnLoad();</script>
</head>
<body>

`)
	tmpl2 = []byte(`

</body>
</html>
`)
)

func run(args []string) int {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stdout, "usage: %s <markdown file>\n", name)
		flag.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitOK
		}
		return exitNG
	}
	args = fs.Args()

	if len(args) != 1 {
		fs.Usage()
		return exitNG
	}
	filename := args[0]

	if filepath.Ext(filename) != ".md" && filepath.Ext(filename) != ".markdown" {
		fs.Usage()
		return exitNG
	}

	wd, _ := os.Getwd()
	b, err := ioutil.ReadFile(filepath.Join(wd, args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}

	// run gom2h
	out, err := gom2h.Run(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}

	// output html
	outFile, err := os.Create(filepath.Base(filename[:len(filename)-len(filepath.Ext(filename))]) + ".html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}
	writer := bufio.NewWriter(outFile)
	if _, err = writer.Write(tmpl1); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}
	if _, err = writer.Write(out); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}
	if _, err = writer.Write(tmpl2); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}
	writer.Flush()

	return exitOK
}
