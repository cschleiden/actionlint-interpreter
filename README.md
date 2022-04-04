# actionlint-interpreter

Simple expression interpreter for _GitHub Actions Expressions_ operating on the AST produced by https://github.com/rhysd/actionlint.

## Usage

```golang
expression := "input.foo <= input.bar"

// Lex & Parse
lexer := actionlint.NewExprLexer(expression + "}}")
parser := actionlint.NewExprParser()
n, perr := parser.Parse(lexer)
if perr != nil {
  panic(perr)
}

// Evaluate expressions
result, err := Evaluate(n, ContextData{
  "input": ContextData{
    "foo": float64(1),
    "bar": float64(2),
  },
})
if err != nil {
  panic(err)
}

fmt.Println(result.Value)
// Output: true
```

### TODO

Not everything is implemented yet:

#### Context access

- [x] Finish object & array access
- [ ] Wildcard access (`inputs.*.foo`)

#### Functions

- [ ] contains
- [x] startsWith
- [x] endsWith
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