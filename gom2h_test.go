package gom2h

import (
	"bytes"
	"testing"
)

func TestEmphasis(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`*em*`, []byte(`<p><em>em</em></p>`)},
		{`This is *em* sample1.`, []byte(`<p>This is <em>em</em> sample1.</p>`)},
		{`This is *multiple* *em* sample2.`, []byte(`<p>This is <em>multiple</em> <em>em</em> sample2.</p>`)},
	}

	for _, tt := range testcases {
		actual, err := Run([]byte(tt.input))
		if err != nil {
			t.Errorf("unexpected err: %v\n", err)
		}
		if !bytes.Equal(tt.expected, actual) {
			t.Errorf("expected %v, but got %v\n", string(tt.expected), string(actual))
		}
	}
}

func TestStrong(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`**strong**`, []byte(`<p><strong>strong</strong></p>`)},
		{`This is **strong** sample1.`, []byte(`<p>This is <strong>strong</strong> sample1.</p>`)},
		{`This is **multiple** **strong** sample2.`, []byte(`<p>This is <strong>multiple</strong> <strong>strong</strong> sample2.</p>`)},
	}

	for _, tt := range testcases {
		actual, err := Run([]byte(tt.input))
		if err != nil {
			t.Errorf("unexpected err: %v\n", err)
		}
		if !bytes.Equal(tt.expected, actual) {
			t.Errorf("expected %v, but got %v\n", string(tt.expected), string(actual))
		}
	}
}

func TestHeader(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`# Header1`, []byte(`<h1>Header1</h1>`)},
		{`## Header2`, []byte(`<h2>Header2</h2>`)},
		{`### Header3`, []byte(`<h3>Header3</h3>`)},
		{`#### Header4`, []byte(`<h4>Header4</h4>`)},
		{`##### Header5`, []byte(`<h5>Header5</h5>`)},
		{`###### Header6`, []byte(`<h6>Header6</h6>`)},
		{`####### Header7`, []byte(`<p>####### Header7</p>`)}, // no header tag
		{`# *em* header`, []byte(`<h1><em>em</em> header</h1>`)},
	}

	for _, tt := range testcases {
		actual, err := Run([]byte(tt.input))
		if err != nil {
			t.Errorf("unexpected err: %v\n", err)
		}
		if !bytes.Equal(tt.expected, actual) {
			t.Errorf("expected %v, but got %v\n", string(tt.expected), string(actual))
		}
	}
}

func TestBlockquote(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`> quote level1`, []byte(`<blockquote><p>quote level1</p></blockquote>`)},
		{`>> quote level2`, []byte(`<blockquote><blockquote><p>quote level2</p></blockquote></blockquote>`)},
		{`> *em* quote`, []byte(`<blockquote><p><em>em</em> quote</p></blockquote>`)},
	}

	for _, tt := range testcases {
		actual, err := Run([]byte(tt.input))
		if err != nil {
			t.Errorf("unexpected err: %v\n", err)
		}
		if !bytes.Equal(tt.expected, actual) {
			t.Errorf("expected %v, but got %v\n", string(tt.expected), string(actual))
		}
	}
}
