# actionlint-interpreter

Simple expression interpreter for _GitHub Actions Expressions_ operating on the AST produced by https://github.com/rhysd/actionlint.


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