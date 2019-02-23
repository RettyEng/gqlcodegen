package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"

	"github.com/RettyInc/gqlcodegen/parser2"
)

func main() {
	p := parser2.NewParser(os.Stdin)
	ts := p.Parse()
	spew.Config.Indent = "    "
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableMethods = true
	spew.Dump(ts)
}
