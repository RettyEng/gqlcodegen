package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"

	"github.com/RettyEng/gqlcodegen/gql"
)

type Package struct {
	Name string
	Path string
}

type Config struct {
	TypeSystem        *gql.TypeSystem
	EnumPackagePrefix string
	ScalarPackage     string
	Package           *Package
}

type Generator struct {
	config *Config
	buff   *bytes.Buffer
}

func NewGenerator(srcInfo *Config) *Generator {
	return &Generator{
		config: srcInfo,
		buff:   bytes.NewBuffer(nil),
	}
}

func (g *Generator) Config() *Config {
	return g.config
}

func (g *Generator) GenerateSource(syntax interface{}) {
	switch def := syntax.(type) {
	case *gql.Enum:
		generateEnum(g, def)
	case *gql.Object:
		generateType(g, def)
	default:
		log.Fatalf("unsupported value %v", def)
	}
}

func (g *Generator) Printf(fmtStr string, args ...interface{}) {
	_, e := fmt.Fprintf(g.buff, fmtStr, args...)
	if e != nil {
		log.Fatal(e)
	}
}

func (g *Generator) Println(args ...interface{}) {
	_, e := fmt.Fprintln(g.buff, args...)
	if e != nil {
		log.Fatal(e)
	}
}

func (g *Generator) Format() {
	src, e := format.Source(g.buff.Bytes())
	if e != nil {
		log.Fatal(e)
	}
	g.buff = bytes.NewBuffer(src)
}

func (g *Generator) ClearBuff() {
	g.buff = bytes.NewBuffer(nil)
}

func (g *Generator) WriteToFile(path string) {
	f, e := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModeAppend|0644)
	if e != nil {
		log.Fatal(e)
	}
	_, e = f.Write(g.buff.Bytes())
	if e != nil {
		log.Fatal(e)
	}
}
