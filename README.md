# actionlint-interpreter

Simple expression interpreter for _GitHub Actions Expressions_ operating on the AST produced by https://github.com/rhysd/actionlint.

## Usage

```golang
expression := "1 <= 2"

// Lex & Parse
lexer := actionlint.NewExprLexer(expression + "}}")
parser := actionlint.NewExprParser()
n, perr := parser.Parse(lexer)
if perr != nil {
  panic(perr)
}

// Evaluate expressions
result, err := Evaluate(n, map[string]interface{}{})
if err != nil {
  panic(err)
}

fmt.Println(result.Value)
// Output: true
```

### TODO

Not everything is implemented yet:

#### Context access

- [ ] Finish object & array access
- [ ] Wildcard access (`inputs.*.foo`)

#### Functions

- [ ] contains
- [x] startsWith
- [ ] endsWith
- [ ] format
- [x] join
- [ ] toJSON
- [ ] fromJSON
- [ ] hashFiles

Status check functions:

- [ ] success
- [ ] always
- [ ] cancelled
- [ ] failure