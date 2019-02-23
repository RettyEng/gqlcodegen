package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type ImplementsExpression interface {
	Eval(object *gql.Object)
}

type DefineFieldExpression interface {
	Eval(object *gql.Object)
}
