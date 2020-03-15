# gom2h

convert markdown to html

## Usage

As a cli tool

```
$ go get -u github.com/matsuyoshi30/gom2h/cmd/gom2h # Installation

$ gom2h <path/to/markdownfile>

$ gom2h -css <path/to/cssfile> <path/to/markdownfile> # specify css
```

[default css](https://github.com/sindresorhus/github-markdown-css)

---

As a go library

```
$ go get -u github.com/matsuyoshi30/gom2h # Installation
```

```go
import (
        "fmt"

        "github.com/matsuyoshi30/gom2h"
)

func main() {
        output, err := gom2h.Run([]byte(`### Header3 with *em* and **strong**`))
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println(string(output)) // -> <h3>Header3 with <em>em</em> and <strong>strong</strong></h3>
}
```

## Support

- [x] Header
- [x] Paragraph
- [x] Emphasis
- [x] Strong
- [x] Link
- [x] List (Unorder)
- [x] Code Block
  - [x] Syntax highlight (only when converting file)

## License

MIT
