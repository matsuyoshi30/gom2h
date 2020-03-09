package gom2h

import (
	"errors"
	"fmt"
	"regexp"
)

// main entry point
func Run(input []byte) ([]byte, error) {
	l, err := conv(input)
	if err != nil {
		return nil, err
	}

	return render(l), nil
}

// convert markdown line to Line

type LineType int

const (
	Header LineType = iota
	Paragraph
	NewLine
)

type Line struct {
	ty  LineType
	lv  int
	val []byte
	dep int
}

var (
	headerExp = regexp.MustCompile(`^(#){1,6} (.+)`)
)

func conv(line []byte) (Line, error) {
	// block

	// inline
	if headerExp.Match(line) {
		loc := headerExp.FindSubmatchIndex(line)
		// ## Header2
		// -> line[loc[0]:loc[3]] // ##
		// -> line[loc[4]:loc[5]] // Header2
		return Line{Header, loc[3], line[loc[4]:loc[5]], 0}, nil
	}

	return Line{}, errors.New("unknown line")
}

// render html from Line

type TagType int

const (
	TagHeader TagType = iota
)

func render(line Line) []byte {
	// render html
	if line.ty == Header {
		return []byte(fmt.Sprintf(`<h%d>%s</h%d>`, line.lv, line.val, line.lv))
	}

	return nil
}
