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
		{`This is _other_ em.`, []byte(`<p>This is <em>other</em> em.</p>`)},
		{`This is _not* em.`, []byte(`<p>This is _not* em.</p>`)},
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
		{`This is **other** strong.`, []byte(`<p>This is <strong>other</strong> strong.</p>`)},
		{`This is **not__ strong.`, []byte(`<p>This is **not__ strong.</p>`)},
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

func TestEmphasisAndStrong(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`***emphasis and strong***`, []byte(`<p><em><strong>emphasis and strong</strong></em></p>`)},
		{`This is ***emphasis and strong*** sample1.`, []byte(`<p>This is <em><strong>emphasis and strong</strong></em> sample1.</p>`)},
		{`This is ***multiple*** ***emphasis and strong*** sample2.`, []byte(`<p>This is <em><strong>multiple</strong></em> <em><strong>emphasis and strong</strong></em> sample2.</p>`)},
		{`___not***`, []byte(`<p>___not***</p>`)},
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

func TestLink(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`![image](/path/to/image)`, []byte(`<p><img src="/path/to/image" alt="image" /></p>`)},
		{`[link](https://example.org/)`, []byte(`<p><a href="https://example.org/">link</a></p>`)},
		{`This is [link](https://example.org/) test.`, []byte(`<p>This is <a href="https://example.org/">link</a> test.</p>`)},
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

func TestList(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{`- list1`, []byte(`<ul>
<li>list1</li>
</ul>`)},
		{`- list1
- list2`, []byte(`<ul>
<li>list1</li>
<li>list2</li>
</ul>`)},
		{`- list1
- list2
  - list2-1
- list3`, []byte(`<ul>
<li>list1</li>
<li>list2</li>
<ul>
<li>list2-1</li>
</ul>
<li>list3</li>
</ul>`)},
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

func TestCodeSpan(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{"`cs sample`", []byte(`<p><code>cs sample</code></p>`)},
		{"This is `cs sample` sentence.", []byte(`<p>This is <code>cs sample</code> sentence.</p>`)},
		{"This is `__emphasis in codespan__` sentence.", []byte(`<p>This is <code>__emphasis in codespan__</code> sentence.</p>`)},
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

func TestCodeFence(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{"```" + `
code fence
` + "```", []byte(`<pre><code>code fence
</code></pre>`)},
		{"```go" + `
fmt.Println("Hello world")
` + "```", []byte(`<pre><code class="go">fmt.Println("Hello world")
</code></pre>`)},
		{"```" + `
- List1
- List2
` + "```", []byte(`<pre><code>- List1
- List2
</code></pre>`)},
		{"```" + `
## header2
#### header4
` + "```", []byte(`<pre><code>## header2
#### header4
</code></pre>`)},
		{"```" + `
> blockquote
>> blockquote2
` + "```", []byte(`<pre><code>&gt; blockquote
&gt;&gt; blockquote2
</code></pre>`)},
		{"```" + `
__strong__
` + "```", []byte(`<pre><code>__strong__
</code></pre>`)},
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

func TestHTMLBlock(t *testing.T) {
	testcases := []struct {
		input    string
		expected []byte
	}{
		{
			`<blockquote>June 19, 2021</blockquote>`,
			[]byte(`<blockquote>June 19, 2021</blockquote>`),
		},
		{
			`<blockquote>_June 19, 2021_</blockquote>`,
			[]byte(`<blockquote>_June 19, 2021_</blockquote>`),
		},
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
