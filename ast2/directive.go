package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineDirectiveArgsExpression interface {
	Eval(*gql.Directive)
}
