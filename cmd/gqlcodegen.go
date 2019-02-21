package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/RettyInc/gqlcodegen/lexer"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	t := lexer.NewLexer(r)
	for _, c := range t.Tokenize() {
		fmt.Printf("%10s: %s\n", c.Type, c.Value)
	}
}
