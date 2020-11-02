package gom2h

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

var nl = []byte("\n")

// main entry point
func Run(input []byte) ([]byte, error) {
	input = bytes.TrimSpace(input)

	conved := make([]Line, 0)
	for _, line := range bytes.Split(input, nl) {
		if len(line) == 0 {
			continue
		}
		l, err := conv(line)
		if err != nil {
			return nil, err
		}
		conved = append(conved, l)
	}

	return render(conved), nil
}

// convert markdown line to Line

type LineType int

const (
	Header LineType = iota
	Blockquote
	List
	CodeFence
	Paragraph
	NewLine
)

type Line struct {
	ty  LineType
	lv  int
	val []byte
	dep int
}

// https://play.golang.org/p/igrR4P6blOD

var (
	headerExp     = regexp.MustCompile(`^(#){1,6} (.+)`)
	blockquoteExp = regexp.MustCompile(`^(>+)(.+)`)
	emphasisExp   = regexp.MustCompile(`.*([\*]([^\*]+)[\*]).*|.*([_]([^_]+)[_]).*`)
	strongExp     = regexp.MustCompile(`.*([\*]{2}([^\*]+)[\*]{2}).*|.*([_]{2}([^_]+)[_]{2}).*`)
	imageExp      = regexp.MustCompile(`^!.*(\[.+\])(\(.+\)).*`)
	linkExp       = regexp.MustCompile(`.*(\[.+\])(\(.+\)).*`)
	listExp       = regexp.MustCompile(`^ *(- )(.+)`)
	codespanExp   = regexp.MustCompile("[^`]*`([^`]+)`[^`]*")
	codefenceExp  = regexp.MustCompile("^```(.*)")
)

func conv(line []byte) (Line, error) {
	inCodeSpan := false

	// inline
	for codespanExp.Match(line) {
		loc := codespanExp.FindSubmatchIndex(line)
		// This is `cs sample`.
		// -> line[loc[2]:loc[3]] // cs sample
		line = []byte(fmt.Sprintf(`%s<code>%s</code>%s`, line[:loc[2]-1], line[loc[2]:loc[3]], line[loc[3]+1:]))
		inCodeSpan = true
	}

	if !inCodeSpan {
		for strongExp.Match(line) {
			loc := strongExp.FindSubmatchIndex(line)
			// This is *em* sample
			// -> line[loc[2]:loc[3]] // **st**
			// -> line[loc[4]:loc[5]] // st
			// -> line[loc[6]:loc[7]] // __st__
			// -> line[loc[8]:loc[9]] // em
			s := loc[2]
			c := loc[3]
			ts := loc[4]
			tc := loc[5]
			if s == -1 && c == -1 && ts == -1 && tc == -1 {
				s = loc[6]
				c = loc[7]
				ts = loc[8]
				tc = loc[9]
			}
			bef := []byte(fmt.Sprintf(`%s`, line[loc[0]:s]))
			target := []byte(fmt.Sprintf(`<strong>%s</strong>`, line[ts:tc]))
			aft := []byte(fmt.Sprintf(`%s`, line[c:]))

			line = append(bef, append(target, aft...)...)
		}

		for emphasisExp.Match(line) {
			loc := emphasisExp.FindSubmatchIndex(line)
			// This is *em* sample
			// -> line[loc[2]:loc[3]] // *em*
			// -> line[loc[4]:loc[5]] // em
			// -> line[loc[6]:loc[7]] // _em_
			// -> line[loc[8]:loc[9]] // em
			s := loc[2]
			c := loc[3]
			ts := loc[4]
			tc := loc[5]
			if s == -1 && c == -1 && ts == -1 && tc == -1 {
				s = loc[6]
				c = loc[7]
				ts = loc[8]
				tc = loc[9]
			}
			bef := []byte(fmt.Sprintf(`%s`, line[loc[0]:s]))
			target := []byte(fmt.Sprintf(`<em>%s</em>`, line[ts:tc]))
			aft := []byte(fmt.Sprintf(`%s`, line[c:]))

			line = append(bef, append(target, aft...)...)
		}

		for imageExp.Match(line) {
			loc := linkExp.FindSubmatchIndex(line)
			// ![image](/path/to/image)
			// -> line[loc[2]:loc[3]] // [image]
			// -> line[loc[4]:loc[5]] // (/path/to/image)
			line = []byte(fmt.Sprintf(`<img src="%s" alt="%s" />`, line[loc[4]+1:loc[5]-1], line[loc[2]+1:loc[3]-1]))
		}

		for linkExp.Match(line) {
			loc := linkExp.FindSubmatchIndex(line)
			// This is [link](https://example.org/)
			// -> line[loc[2]:loc[3]] // [link]
			// -> line[loc[4]:loc[5]] // (https://example.org/)
			bef := []byte(fmt.Sprintf(`%s`, line[loc[0]:loc[2]]))
			target := []byte(fmt.Sprintf(`<a href="%s">%s</a>`, line[loc[4]+1:loc[5]-1], line[loc[2]+1:loc[3]-1]))
			aft := []byte(fmt.Sprintf(`%s`, line[loc[5]:]))

			line = append(bef, append(target, aft...)...)
		}
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

	if listExp.Match(line) {
		loc := listExp.FindSubmatchIndex(line)
		// - list
		// -> line[loc[4]:loc[5]] // list
		return Line{List, 0, line[loc[4]:loc[5]], loc[2] / 2}, nil
	}

	if codefenceExp.Match(line) {
		loc := codefenceExp.FindSubmatchIndex(line)
		// ```go
		// -> line[loc[2]:loc[3]] // go
		var lang []byte
		if loc[2] != loc[3] {
			lang = line[loc[2]:loc[3]]
		}
		return Line{CodeFence, 0, lang, 0}, nil
	}

	return Line{Paragraph, 0, line, 0}, nil
}

// render html from Line

func render(lines []Line) []byte {
	ret := make([]byte, 0)

	inCodeFence := false
	for idx, line := range lines {
		// render html
		switch line.ty {
		case Header:
			if inCodeFence {
				ret = append(ret, []byte(fmt.Sprintf(`%s %s`, strings.Repeat("#", line.lv), line.val))...)
			} else {
				ret = append(ret, []byte(fmt.Sprintf(`<h%d>%s</h%d>`, line.lv, line.val, line.lv))...)
			}

		case Blockquote:
			if inCodeFence {
				ret = append(ret, []byte(fmt.Sprintf(`%s%s`, strings.Repeat("&gt;", line.lv), line.val))...)
			} else {
				var stag string
				var ctag string
				for i := 0; i < line.lv; i++ {
					stag = fmt.Sprintf(`<blockquote>%s`, stag)
					ctag = fmt.Sprintf(`%s</blockquote>`, ctag)
				}
				ret = append(ret, []byte(fmt.Sprintf(`%s<p>%s</p>%s`, stag, bytes.TrimSpace(line.val), ctag))...)
			}

		case List:
			if inCodeFence {
				ret = append(ret, []byte(fmt.Sprintf(`- %s`, line.val))...)
			} else {
				if (idx > 0 && lines[idx-1].ty != List) || idx == 0 {
					ret = append(ret, []byte(`<ul>`)...)
					ret = newline(ret)
				}

				if idx > 0 && lines[idx-1].dep < line.dep {
					ret = append(ret, []byte(`<ul>`)...)
					ret = newline(ret)
				}

				ret = append(ret, []byte(fmt.Sprintf(`<li>%s</li>`, line.val))...)

				if idx < len(lines)-1 && line.dep > lines[idx+1].dep {
					for d := line.dep - lines[idx+1].dep; d > 0; d-- {
						ret = newline(ret)
						ret = append(ret, []byte(`</ul>`)...)
					}
				}
				if idx == len(lines)-1 && line.dep > 0 {
					for d := line.dep; d > 0; d-- {
						ret = newline(ret)
						ret = append(ret, []byte(`</ul>`)...)
					}
				}

				if (idx < len(lines)-1 && lines[idx+1].ty != List) || idx == len(lines)-1 {
					ret = newline(ret)
					ret = append(ret, []byte(`</ul>`)...)
				}
			}

		case CodeFence:
			if !inCodeFence {
				if line.val != nil {
					ret = append(ret, []byte(fmt.Sprintf(`<pre><code class="%s">`, line.val))...)
				} else {
					ret = append(ret, []byte(`<pre><code>`)...)
				}
			} else {
				ret = append(ret, []byte(`</code></pre>`)...)
			}
			inCodeFence = !inCodeFence
			continue

		case Paragraph:
			if !inCodeFence {
				ret = append(ret, []byte(fmt.Sprintf(`<p>%s</p>`, line.val))...)
			} else {
				ret = append(ret, []byte(fmt.Sprintf(`%s`, line.val))...)
			}
		}

		if idx != len(lines)-1 {
			ret = newline(ret)
		}
	}

	return ret
}

func newline(ret []byte) []byte {
	return append(ret, nl...)
}
