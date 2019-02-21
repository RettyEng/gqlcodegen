package main

import (
	"bufio"
	"fmt"
	"github.com/RettyInc/gqlcodegen/lexer"
	"os"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	t := lexer.NewLexer(r)
	for _, c := range t.Tokenize() {
		if c == "\n" {
			c = "\\n"
		}
		fmt.Printf("%v\t%x\n", c, c)
	}
}
