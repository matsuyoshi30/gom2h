package gom2h

import (
	"bytes"
	"testing"
)

func TestHeader(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
		isErr    bool
	}{
		{`# Header1`, []byte(`<h1>Header1</h1>`), false},
		{`## Header2`, []byte(`<h2>Header2</h2>`), false},
		{`### Header3`, []byte(`<h3>Header3</h3>`), false},
		{`#### Header4`, []byte(`<h4>Header4</h4>`), false},
		{`##### Header5`, []byte(`<h5>Header5</h5>`), false},
		{`###### Header6`, []byte(`<h6>Header6</h6>`), false},
		{`####### Header7`, nil, true},
	}

	for _, tt := range testcases {
		actual, err := Run([]byte(tt.input))
		if !bytes.Equal(tt.expected, actual) {
			t.Errorf("expected %v, but got %v\n", string(tt.expected), string(actual))
		}
		if (tt.isErr && err == nil) || (!tt.isErr && err != nil) {
			t.Errorf("expected %v, but got %v\n", string(tt.expected), string(actual))
		}
	}
}
