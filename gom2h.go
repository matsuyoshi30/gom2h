package gom2h

import (
	"bytes"
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
	Blockquote
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
	headerExp     = regexp.MustCompile(`^(#){1,6} (.+)`)
	blockquoteExp = regexp.MustCompile(`^(>+)(.+)`)
	emphasisExp   = regexp.MustCompile(`.*([\*_]([^\*_]+)[\*_]).*`)
)

func conv(line []byte) (Line, error) {
	// inline
	for emphasisExp.Match(line) {
		loc := emphasisExp.FindSubmatchIndex(line)
		// This is *em* sample
		// -> line[loc[2]:loc[3]] // *em*
		// -> line[loc[4]:loc[5]] // em
		bef := []byte(fmt.Sprintf(`%s`, line[loc[0]:loc[2]]))
		target := []byte(fmt.Sprintf(`<em>%s</em>`, line[loc[4]:loc[5]]))
		aft := []byte(fmt.Sprintf(`%s`, line[loc[3]:]))

		line = append(bef, append(target, aft...)...)
	}

	// block
	if headerExp.Match(line) {
		loc := headerExp.FindSubmatchIndex(line)
		// ## Header2
		// -> line[loc[0]:loc[3]] // ##
		// -> line[loc[4]:loc[5]] // Header2
		return Line{Header, loc[3], line[loc[4]:loc[5]], 0}, nil
	}

	if blockquoteExp.Match(line) {
		loc := blockquoteExp.FindSubmatchIndex(line)
		// > quote
		// -> line[loc[0]:loc[3]] // >
		// -> line[loc[4]:loc[4]] // quote
		return Line{Blockquote, loc[3], line[loc[4]:loc[5]], 0}, nil
	}

	return Line{Paragraph, 0, line, 0}, nil
}

// render html from Line

type TagType int

func render(line Line) []byte {
	// render html
	if line.ty == Header {
		return []byte(fmt.Sprintf(`<h%d>%s</h%d>`, line.lv, line.val, line.lv))
	}
	if line.ty == Blockquote {
		var stag string
		var ctag string
		for i := 0; i < line.lv; i++ {
			stag = fmt.Sprintf(`<blockquote>%s`, stag)
			ctag = fmt.Sprintf(`%s</blockquote>`, ctag)
		}
		return []byte(fmt.Sprintf(`%s<p>%s</p>%s`, stag, bytes.TrimSpace(line.val), ctag))
	}
	if line.ty == Paragraph {
		return []byte(fmt.Sprintf(`<p>%s</p>`, line.val))
	}

	return nil
}
