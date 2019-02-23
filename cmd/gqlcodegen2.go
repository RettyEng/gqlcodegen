package main

import (
	"fmt"
	"os"

	"github.com/RettyInc/gqlcodegen/lexer"
)

func main() {
	l := lexer.NewLexer(os.Stdin)

	line := 1

	for t := l.Next(); t != nil; t = l.Next() {
		li, co := t.LineCol()
		if line != li {
			fmt.Println()
			line = li
		}
		fmt.Printf("%s(%s)[%d:%d] ", t.Type().String(), t.Value(), li, co)
	}
}
