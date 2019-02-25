package ast

import (
	"github.com/RettyInc/gqlcodegen/ast/directive"

	"github.com/RettyInc/gqlcodegen/gql"
)

type DirectiveInternalExpression interface {
	Eval(*gql.Directive)
}

type DefineDirectiveLocationExpression struct {
	Location directive.Location
}

func (d *DefineDirectiveLocationExpression) Eval(direc *gql.Directive) {
	direc.Location = append(direc.Location, d.Location)
}
