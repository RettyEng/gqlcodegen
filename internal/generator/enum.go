package generator

import (
	"strings"

	"github.com/RettyInc/gqlcodegen/ast"
)

func generateEnum(g *Generator, def *ast.EnumDef) {
	g.Printf(commentOnTop)
	generateEnumPackageSection(g, def)
	g.Println()
	generateEnumTypeDefSection(g, def)
	g.Println()
	generateEnumConstSection(g, def)
	g.Println()
	g.Printf("//go:generate stringer -type=%s\n", capitalizeFirst(def.Name()))
}

func generateEnumPackageSection(g *Generator, def *ast.EnumDef) {
	g.Printf("package %s\n", strings.ToLower(def.Name()))
}

func generateEnumTypeDefSection(g *Generator, def *ast.EnumDef) {
	g.Printf("type %s int\n", capitalizeFirst(def.Name()))
}

func generateEnumConstSection(g *Generator, def *ast.EnumDef) {
	g.Println("const(")
	generateEnumConstBody(g, def)
	g.Println(")")
}

func generateEnumConstBody(g *Generator, def *ast.EnumDef) {
	entries := def.Entries()
	for i, e := range entries {
		g.Printf("%s", capitalizeFirst(e.Name()))
		if i == 0 {
			g.Printf(" %s = iota", capitalizeFirst(def.Name()))
		}
		g.Println()
	}
}
