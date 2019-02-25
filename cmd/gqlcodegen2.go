package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/RettyInc/gqlcodegen/parser"
)

func main() {
	p := parser.NewParser(os.Stdin)
	ts := p.ParseSchema()
	spew.Config.Indent = "    "
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableMethods = true
	spew.Dump(ts)
}
