package generator

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/RettyInc/gqlcodegen/gql"
)

func generateType(g *Generator, def *gql.Object) {
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

func generateArgStructs(g *Generator, def *gql.Object) {
	for _, f := range def.Fields {
		if len(f.Args) > 0 {
			generateArgStruct(g, def, f)
			g.Println()
		}
	}
}

func generateArgStruct(g *Generator, t *gql.Object, f *gql.ObjectField) {
	g.Printf("type %s struct {\n", argStructName(f, t))
	for _, a := range f.Args {
		generateComment(g, a)
		g.Printf("%s %s\n", capitalizeFirst(a.Name), refToString(g, a.Type))
	}
	g.Println("}")
}

func generateTypeDefinition(g *Generator, def *gql.Object) {
	generateComment(g, def)
	g.Printf("type %s interface {\n", convertResolverName(def.Name))
	for _, f := range def.Fields {
		generateField(g, f, def)
	}
	g.Println("}")
}

func generateField(g *Generator, f *gql.ObjectField, t *gql.Object) {
	name := capitalizeFirst(f.Name)
	if f.Type.IsNullable {
		g.Println()
		generateComment(g, f)
		g.Printf("// Return value of %s is nullable\n", name)
	} else {
		generateComment(g, f)
	}
	g.Printf("%s(", name)
	if len(f.Args) > 0 {
		g.Printf(
			"context.Context, %s",
			argStructName(f, t),
		)
	}
	g.Printf(") %s\n", refToString(g, f.Type))
}

func argStructName(f *gql.ObjectField, t *gql.Object) string {
	return convertResolverName(t.Name) + "_" + capitalizeFirst(f.Name) + "_Arg"
}

func convertResolverName(name string) string {
	return capitalizeFirst(name) + "Resolver"
}

func refToString(g *Generator, ref *gql.TypeRef) string {
	n := ref.Name
	if _, ok := g.Config().TypeSystem.ObjectTypes[n]; ok {
		return convertResolverName(ref.Name)
	}
	if _, ok := g.Config().TypeSystem.EnumTypes[n]; ok {
		n = strings.ToLower(n) + "." + capitalizeFirst(n)
		if ref.IsNullable {
			n = "*" + n
		}
		return n
	}
	if _, ok := g.Config().TypeSystem.ScalarTypes[n]; ok {
		n = path.Base(g.Config().ScalarPackage) + "." + capitalizeFirst(n)
		if ref.IsNullable {
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
		n = "[]" + refToString(g, ref.InnerType)
	default:
		log.Fatalf("unknown type %s", n)
	}
	if ref.IsNullable {
		n = "*" + n
	}
	return n
}

func generateResolverPackageSection(g *Generator) {
	g.Printf("package %s\n", g.Config().Package.Name)
}

func findTypeName(ref []*gql.TypeRef) []string {
	var ret []string
	for _, r := range ref {
		ret = append(ret, r.Name)
		if r.InnerType == nil {
			continue
		}
		ret = append(ret, findTypeName([]*gql.TypeRef{r.InnerType})...)
	}
	return ret
}

func typeNameMap(def *gql.Object) map[string]struct{} {
	var refs []*gql.TypeRef
	for _, f := range def.Fields {
		refs = append(refs, f.Type)
		for _, a := range f.Args {
			refs = append(refs, a.Type)
		}
	}
	types := map[string]struct{}{}
	for _, t := range findTypeName(refs) {
		types[t] = struct{}{}
	}
	return types
}

func generateImportSection(g *Generator, def *gql.Object) {
	typesMap := typeNameMap(def)
	needContext := false
	for _, f := range def.Fields {
		if len(f.Args) > 0 {
			needContext = true
			break
		}
	}
	var imported []string
	if needContext {
		imported = append(imported, "\"context\"")
	}
	for k := range g.Config().TypeSystem.EnumTypes {
		if _, ok := typesMap[k]; ok {
			imported = append(
				imported,
				fmt.Sprintf("\"%s\"", path.Join(g.Config().EnumPackagePrefix, strings.ToLower(k))),
			)
		}
	}
	for k := range g.Config().TypeSystem.ScalarTypes {
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
