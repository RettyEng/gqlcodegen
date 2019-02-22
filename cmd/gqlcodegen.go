package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/RettyInc/gqlcodegen/ast"
	"github.com/RettyInc/gqlcodegen/internal/generator"
	"github.com/RettyInc/gqlcodegen/parser"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"path"
	"strings"
)

var (
	enumPackagePrefix = flag.String("enum-pkg-prefix", "", "")
	generateTarget = flag.String("target", "", "comma separated")
	schema = flag.String("schema", "", "comma separated")
)

func main() {
	flag.Parse()
	fmt.Println(*enumPackagePrefix)
	fmt.Println(*generateTarget)
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Tests: false,
	}
	args := flag.Args()
	if len(args) == 0 {
		args = []string{
			".",
		}
	} else {
		args = []string{
			path.Dir(args[0]),
		}
	}
	pkgs, _ := packages.Load(cfg, args...)
	pkg := pkgs[0]

	ast := loadAst(*schema)

	scalar := map[string]struct{}{}
	for _, s := range ast.Scalars() {
		scalar[s.Name()] = struct{}{}
	}

	enum := map[string]struct{}{}
	for _, e := range ast.Enums() {
		enum[e.Name()] = struct{}{}
	}

	resolver := map[string]struct{}{}
	for _, r := range ast.Types() {
		resolver[r.Name()] = struct{}{}
	}

	genCfg := &generator.Config{
		ScalarTypes: scalar,
		ResolverTypes: resolver,
		EnumTypes:enum,
		EnumPackagePrefix:"",
		ScalarPackagePrefix:"",
		Package:&generator.Package{
			Name:pkg.Name,
		},
	}

	gen := generator.NewGenerator(genCfg)
	for _, e := range ast.Enums() {
		writeEnum(gen, e)
		gen.ClearBuff()
	}

}

func writeEnum(g *generator.Generator,def *ast.EnumDef) {
	g.GenerateSource(def)
	g.Format()
	dirName := "./" + strings.ToLower(def.Name())
	_ = os.Mkdir(dirName, 0755)
	g.WriteToFile(dirName + "/" + strings.ToLower(def.Name()) + ".go")
}

func loadAst(schemaPath string) *ast.Root {
	f, e := os.Open(schemaPath)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()
	return parser.NewParser(bufio.NewReader(f)).Parse()
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
