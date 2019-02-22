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
	scalarPackage = flag.String("scalar-pkg", "", "")
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
	fmt.Printf("%v", resolver)

	genCfg := &generator.Config{
		ScalarTypes: scalar,
		ResolverTypes: resolver,
		EnumTypes:enum,
		EnumPackagePrefix: *enumPackagePrefix,
		ScalarPackage: *scalarPackage,
		Package:&generator.Package{
			Name:pkg.Name,
		},
	}

	gen := generator.NewGenerator(genCfg)
	for _, e := range ast.Types() {
		writeType(gen, e)
	}
}

func writeEnum(g *generator.Generator,def *ast.EnumDef) {
	g.GenerateSource(def)
	defer g.ClearBuff()
	g.Format()
	dirName := "./" + strings.ToLower(def.Name())
	_ = os.Mkdir(dirName, 0755)
	g.WriteToFile(dirName + "/" + strings.ToLower(def.Name()) + ".go")
}

func writeType(g *generator.Generator, def *ast.TypeDef) {
	g.GenerateSource(def)
	defer g.ClearBuff()
	g.Format()
	g.WriteToFile(strings.ToLower(def.Name()) + ".go")
}

func loadAst(schemaPath string) *ast.Root {
	f, e := os.Open(schemaPath)
	if e != nil {
		log.Fatal(e)
	}
	defer f.Close()
	return parser.NewParser(bufio.NewReader(f)).Parse()
}
