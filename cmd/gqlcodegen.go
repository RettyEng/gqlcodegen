package main

import (
	"bufio"
	"fmt"
	"github.com/RettyInc/gqlcodegen/ast"
	"os"

	"github.com/RettyInc/gqlcodegen/lexer"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	t := lexer.NewLexer(r)
	for _, c := range t.Tokenize() {
		fmt.Printf("%10s: %s\n", c.Type, c.Value)
	}
	root := ast.NewRoot([]ast.Ast{
		ast.NewTypeDef("hoge", nil),
		ast.NewTypeDef("fuga", nil),
		ast.NewEnumDef("foo", nil),
		ast.NewEnumDef("bar", nil),
		ast.NewTypeDef("piyo", nil),
		ast.NewEnumDef("baz", nil),
	})
	for _, t := range root.Types() {
		fmt.Printf("%8s, %s\n", t.AstType().String(), t.Name())
	}
	for _, t := range root.Enums() {
		fmt.Printf("%8s, %s\n", t.AstType().String(), t.Name())
	}
}
