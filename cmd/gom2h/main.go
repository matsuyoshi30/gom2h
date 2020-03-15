package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
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

type Page struct {
	Stylesheet template.CSS
	Content    template.HTML
}

func run(args []string) int {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintf(os.Stdout, "usage: %s <markdown file>\n", name)
		flag.PrintDefaults()
	}

	var cssfile string
	fs.StringVar(&cssfile, "css", "", "path to css file")
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

	// read css
	var style []byte
	if cssfile != "" {
		style, err = ioutil.ReadFile(filepath.Join(wd, cssfile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read css: %v\n", err)
			return exitNG
		}
	} else {
		style = css()
	}

	// run gom2h
	out, err := gom2h.Run(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}

	page := Page{Stylesheet: template.CSS(style), Content: template.HTML(out)}

	tmpl, err := template.New("index").Parse(index)
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
	if err = tmpl.Execute(writer, page); err != nil {
		fmt.Fprintf(os.Stderr, "unexpected error: %v\n", err)
		return exitNG
	}
	writer.Flush()

	return exitOK
}
