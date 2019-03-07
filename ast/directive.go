package ast

import (
	"github.com/RettyEng/gqlcodegen/ast/directive"

	"github.com/RettyEng/gqlcodegen/gql"
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
