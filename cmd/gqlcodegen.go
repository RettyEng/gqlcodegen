package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/RettyInc/gqlcodegen/gql"
	"github.com/RettyInc/gqlcodegen/internal/generator"
	"github.com/RettyInc/gqlcodegen/parser"
)

var (
	fileSuffix        = flag.String("suffix", "_gql", "")
	enumPackagePrefix = flag.String("enum-pkg-prefix", "", "")
	scalarPackage     = flag.String("scalar-pkg", "", "")
	generateTarget    = flag.String("target", "", "comma separated")
	schema            = flag.String("schema", "", "comma separated")
)

func createGenerator(
	packageName, packagePath string, root *gql.TypeSystem,
) *generator.Generator {
	conf := &generator.Config{
		TypeSystem:        root,
		EnumPackagePrefix: *enumPackagePrefix,
		ScalarPackage:     *scalarPackage,
		Package: &generator.Package{
			Name: packageName,
			Path: packagePath,
		},
	}
	return generator.NewGenerator(conf)
}

func main() {
	flag.Parse()

	packagePath, _ := filepath.Abs(".")
	packageName := path.Base(packagePath)

	if args := flag.Args(); len(args) > 0 {
		packagePath = path.Dir(args[0])
		packageName = path.Base(packagePath)
	}
	typeSystem := loadTypeSystem(*schema)
	generator := createGenerator(packageName, packagePath, typeSystem)

	generate(generator)
}

func generate(g *generator.Generator) {
	targets := strings.Split(*generateTarget, ",")
	for _, t := range targets {
		switch t {
		case "enum":
			for _, e := range g.Config().TypeSystem.EnumTypes {
				writeEnum(g, e)
			}
		case "resolver":
			for _, t := range g.Config().TypeSystem.ObjectTypes {
				writeType(g, t)
			}
		default:
			log.Fatalf("unknown target %s", t)
		}
	}
}

func writeEnum(g *generator.Generator, enum *gql.Enum) {
	g.GenerateSource(enum)
	defer g.ClearBuff()
	g.Format()
	dirName := path.Join(g.Config().Package.Path, strings.ToLower(enum.Name))
	_ = os.Mkdir(dirName, 0755)
	g.WriteToFile(
		path.Join(dirName, strings.ToLower(enum.Name)+*fileSuffix+".go"),
	)
}

func writeType(g *generator.Generator, obj *gql.Object) {
	g.GenerateSource(obj)
	defer g.ClearBuff()
	g.Format()
	g.WriteToFile(
		path.Join(g.Config().Package.Path, strings.ToLower(obj.Name)+*fileSuffix+".go"),
	)
}

func loadTypeSystem(schemaPath string) *gql.TypeSystem {
	f, e := os.Open(schemaPath)
	if e != nil {
		log.Fatalf("error occured while loading schema: %v", e)
	}
	defer f.Close()
	return parser.NewParser(bufio.NewReader(f)).ParseAndEvalSchema()
}
