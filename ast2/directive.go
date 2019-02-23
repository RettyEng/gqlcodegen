package ast2

import (
	"github.com/RettyInc/gqlcodegen/ast2/directive"

	"github.com/RettyInc/gqlcodegen/gql"
)

type DirectiveInternalExpression interface {
	Eval(*gql.Directive)
}

type DefineDirectiveLocationExpression struct {
	locations []directive.Location
}

func (d *DefineDirectiveLocationExpression) Eval(direc *gql.Directive) {
	direc.Location = d.locations
}
