package generator

import (
	"fmt"
	"strings"

	"github.com/RettyEng/gqlcodegen/gql"
)

func generateComment(
	g *Generator,
	c gql.Commentable,
) {
	if len(c.GetDescription()) == 0 && len(c.GetDirectives()) == 0 {
		return
	}
	g.Println("\n/*")
	if len(c.GetDescription()) != 0 {
		g.Printf("   Description:\n     %s\n", strings.Replace(c.GetDescription(), "\n", "\n     ", -1))
	}
	if len(c.GetDirectives()) != 0 {
		g.Printf("   Directives:\n")
		for _, d := range c.GetDirectives() {
			g.Printf("     @%s(%s)\n", d.Name, argsStr(d.Args))
		}
	}
	g.Println(" */")
}

func argsStr(args map[string]gql.Value) string {
	var str []string
	for name, value := range args {
		str = append(str, fmt.Sprintf("%s: %s", name, value.Value()))
	}
	return strings.Join(str, ", ")
}
