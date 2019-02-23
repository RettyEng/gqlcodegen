package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineInputObjectFieldExpression interface {
	Eval(object *gql.InputObject)
}
