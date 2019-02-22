package main

import (
	"bufio"
	"flag"
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
	fileSuffix = flag.String("suffix", "_gql", "")
	enumPackagePrefix = flag.String("enum-pkg-prefix", "", "")
	scalarPackage = flag.String("scalar-pkg", "", "")
	generateTarget = flag.String("target", "", "comma separated")
	schema = flag.String("schema", "", "comma separated")
)

func loadPackage(pattern string) *packages.Package {
	cfg := &packages.Config{
		Mode: packages.LoadSyntax,
		Tests: false,
	}
	pkgs, e := packages.Load(cfg, pattern)
	if e != nil {
		log.Fatal(e)
	}
	return pkgs[0]
}

func createGenerator(pkg *packages.Package, root *ast.Root) *generator.Generator {
	conf := &generator.Config{
		EnumTypes:enumerateEnums(root),
		ScalarTypes:enumerateScalars(root),
		ResolverTypes:enumerateResolverTypes(root),
		EnumPackagePrefix:*enumPackagePrefix,
		ScalarPackage:*scalarPackage,
		Package:&generator.Package{
			Name: pkg.Name,
			Path: pkg.PkgPath,
		},
	}
	return generator.NewGenerator(conf)
}

func enumerateScalars(root *ast.Root) map[string]struct{} {
	scalar := map[string]struct{}{}
	for _, s := range root.Scalars() {
		scalar[s.Name()] = struct{}{}
	}
	return scalar
}

func enumerateEnums(root *ast.Root) map[string]struct{} {
	enum := map[string]struct{}{}
	for _, e := range root.Enums() {
		enum[e.Name()] = struct{}{}
	}
	return enum
}

func enumerateResolverTypes(root *ast.Root) map[string]struct{} {
	resolver := map[string]struct{}{}
	for _, r := range root.Types() {
		resolver[r.Name()] = struct{}{}
	}
	return resolver
}

func main() {
	flag.Parse()

	pattern := "."
	if args := flag.Args(); len(args) > 0 {
		pattern = path.Dir(args[0])
	}
	pkg := loadPackage(pattern)
	rootAst := loadAst(*schema)
	generator := createGenerator(pkg, rootAst)

	generate(generator, rootAst)
}

func generate(g *generator.Generator, root *ast.Root) {
	targets := strings.Split(*generateTarget, ",")
	for _, t := range targets {
		switch t {
		case "enum":
			for _, e := range root.Enums() {
				writeEnum(g, e)
			}
		case "resolver":
			for _, t := range root.Types() {
				writeType(g, t)
			}
		default:
			log.Fatal("unknown target %s", t)
		}
	}
}

func writeEnum(g *generator.Generator,def *ast.EnumDef) {
	g.GenerateSource(def)
	defer g.ClearBuff()
	g.Format()
	dirName := path.Join(g.Config().Package.Path, strings.ToLower(def.Name()))
	e := os.Mkdir(dirName, 0755)
	if e != nil {
		log.Fatalf("could not create directory %v", e)
	}
	g.WriteToFile(
		path.Join(dirName, strings.ToLower(def.Name()) + *fileSuffix + ".go"),
	)
}

func writeType(g *generator.Generator, def *ast.TypeDef) {
	g.GenerateSource(def)
	defer g.ClearBuff()
	g.Format()
	g.WriteToFile(
		path.Join(g.Config().Package.Path, strings.ToLower(def.Name()) + ".go"),
	)
}

func loadAst(schemaPath string) *ast.Root {
	f, e := os.Open(schemaPath)
	if e != nil {
		log.Fatalf("error occured while loading schema: %v", e)
	}
	defer f.Close()
	return parser.NewParser(bufio.NewReader(f)).Parse()
}
