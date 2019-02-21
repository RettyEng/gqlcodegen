package main

import (
	"bufio"
	"fmt"
	"github.com/RettyInc/gqlcodegen/ast"
	"github.com/RettyInc/gqlcodegen/parser"
	"os"
	"strings"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	p := parser.NewParser(r)
	root := p.Parse()
	for _, e := range root.Enums() {
		enums[e.Name()] = struct{}{}
	}
	for _, s := range root.Scalars() {
		scalars[s.Name()] = struct{}{}
	}
	for _, t := range root.Types() {
		printTypes(t)
	}
	fmt.Println()
	fmt.Println()
	for _, e := range root.Enums() {
		printEnum(e)
	}
}

var scalars map[string]struct{} = map[string]struct{}{}
var enums map[string]struct{} = map[string]struct{}{}

func printTypes(t *ast.TypeDef) {
	body := ""
	for _,f := range t.Fields() {
		body += "\t"
		n := f.Name()
		n = strings.ToUpper(n[0:1]) + n[1:]
		body += n
		body += "("
		a := f.Args()
		if len(a) > 0 {
			body += "context.Context, " + t.Name() + "_" + n + "_Args"
		}
		body += ") "
		body += convertTypeName(f.Type()) + "\n"
	}
	fmt.Printf(`
type %s interface {
%s
}
`, t.Name(), body)
}

func printEnum(t *ast.EnumDef) {
	body := ""
	iota := true
	for _, e := range t.Entries() {
		body += "\t"
		body += e.Name()
		if iota {
			body += " " + t.Name() + " = iota"
			iota = false
		}
		body += "\n"
	}
	fmt.Printf(`
type %s int
const (
%s
)
`, t.Name(), body)
}

func convertTypeName(t *ast.TypeRef) string {
	n := t.Name()
	if _, ok := scalars[n]; ok {
		n = "scalar." + t.Name()
	} else if _, ok := enums[n]; ok {
		n = "enum." + strings.ToLower(t.Name()) + "." + t.Name()
	}

	switch n {
	case "Int":
		n = "int"
	case "String":
		n = "string"
	case "Boolean":
		n = "bool"
	case "[]":
		n = "[]" + convertTypeName(t.TypeVars()[0])
	}
	if n == t.Name() {
		return n
	}
	p := ""
	if t.IsNullable() {
		p = "*"
	}
	return p + n
}
