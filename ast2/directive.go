package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineDirectiveArgsExpression interface {
	Eval(*gql.Directive)
}

type DefineDirectiveLocationExpression interface {
	Eval(*gql.Directive)
}