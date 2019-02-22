package generator

import (
	"fmt"
	"github.com/RettyInc/gqlcodegen/ast"
	"log"
	"path"
	"strings"
)

func generateType(g *Generator, def *ast.TypeDef) {
	g.Printf(commentOnTop)
	generateResolverPackageSection(g)
	g.Println()
	generateImportSection(g, def)
	g.Println()
	generateTypeDefinition(g, def)
	g.Println()
	g.Println()
	generateArgStructs(g, def)
}

func generateArgStructs(g *Generator, def *ast.TypeDef) {
	for _, f := range def.Fields() {
		if len(f.Args()) > 0 {
			generateArgStruct(g, def, f)
			g.Println()
		}
	}
}

func generateArgStruct(g *Generator, t *ast.TypeDef, f *ast.TypeFieldDef) {
	g.Printf("type %s struct {\n", argStructName(f, t))
	for _, a := range f.Args() {
		g.Printf("%s %s\n", capitalizeFirst(a.Name()), refToString(g, a.Type()))
	}
	g.Println("}")
}

func generateTypeDefinition(g *Generator, def *ast.TypeDef) {
	g.Printf("type %s interface {\n", convertResolverName(def.Name()))
	for _, f := range def.Fields() {
		generateField(g, f, def)
	}
	g.Println("}")
}

func generateField(g *Generator, f *ast.TypeFieldDef, t *ast.TypeDef) {
	name := capitalizeFirst(f.Name())
	if f.Type().IsNullable() {
		g.Println()
		g.Printf("// Return value of %s is nullable\n", name)
	}
	g.Printf("%s(", name)
	if len(f.Args()) > 0 {
		g.Printf(
			"context.Context, %s",
			argStructName(f, t),
		)
	}
	g.Printf(") %s\n", refToString(g, f.Type()))
}

func argStructName(f *ast.TypeFieldDef, t *ast.TypeDef) string {
	return convertResolverName(t.Name()) + "_" + capitalizeFirst(f.Name()) + "_Arg"
}

func convertResolverName(name string) string {
	return capitalizeFirst(name) + "Resolver"
}

func refToString(g *Generator, ref *ast.TypeRef) string {
	n := ref.Name()
	if _, ok := g.Config().ResolverTypes[n] ; ok {
		return convertResolverName(ref.Name())
	}
	if _, ok := g.Config().EnumTypes[n]; ok {
		n = strings.ToLower(n) + "." + capitalizeFirst(n)
		if ref.IsNullable() {
			n = "*" + n
		}
		return n
	}
	if _, ok := g.Config().ScalarTypes[n]; ok {
		n = path.Base(g.Config().ScalarPackage) + "." + capitalizeFirst(n)
		if ref.IsNullable() {
			n = "*" + n
		}
		return n
	}

	switch n {
	case "Int":
		n = "int"
	case "String":
		n = "string"
	case "Boolean":
		n = "bool"
	case "Float":
		n = "float32"
	case "[]":
		n = "[]" + refToString(g, ref.TypeVars()[0])
	default:
		log.Fatalf("unknown type %s", n)
	}
	if ref.IsNullable() {
		n = "*" + n
	}
	return n
}

func generateResolverPackageSection(g *Generator) {
	g.Printf("package %s\n", g.Config().Package.Name)
}

func findTypeName(ref []*ast.TypeRef) []string {
	var ret []string
	for _, r := range ref {
		ret = append(ret, r.Name())
		ret = append(ret, findTypeName(r.TypeVars())...)
	}
	return ret
}

func typeNameMap(def *ast.TypeDef) map[string]struct{} {
	var refs []*ast.TypeRef
	for _, f := range def.Fields() {
		refs = append(refs, f.Type())
		for _, a := range f.Args() {
			refs = append(refs, a.Type())
		}
	}
	types := map[string]struct{}{}
	for _, t := range findTypeName(refs) {
		types[t] = struct{}{}
	}
	return types
}

func generateImportSection(g *Generator, def *ast.TypeDef) {
	typesMap := typeNameMap(def)
	needContext := false
	for _, f := range def.Fields() {
		if len(f.Args()) > 0 {
			needContext = true
			break
		}
	}
	var imported []string
	if needContext {
		imported = append(imported, "\"context\"")
	}
	for k := range g.Config().EnumTypes {
		if _, ok := typesMap[k]; ok {
			imported = append(
				imported,
				fmt.Sprintf("\"%s\"", path.Join(g.Config().EnumPackagePrefix, strings.ToLower(k))),
			)
		}
	}
	for k := range g.Config().ScalarTypes {
		if _, ok := typesMap[k]; ok {
			imported = append(
				imported,
				fmt.Sprintf("\"%s\"", strings.Trim(g.Config().ScalarPackage, "/")),
			)
			break
		}
	}

	if len(imported) > 0 {
		g.Println("import(")
		for _, p := range imported {
			g.Println(p)
			if p == "\"context\"" {
				g.Println()
			}
		}
		g.Println(")")
	}

}
