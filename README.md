# gqlcodegen
This is a code-generator expected to be used with [github.com/graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go).

# Usage
## Installation
Install binary with following command.
```sh
go get -u github.com/RettyEng/gqlcodegen/cmd/gqlcodegen
```

## Generating codes
See [resolver.go](./example/resolver.go) and [enum.go](./example/enum/enum.go)

## Custom generator
If you want to work with another graphql library, you can define custom generator with parser in this package.
[example](./cmd/gqlcodegen/main.go)
