package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineInterfaceFieldExpression interface {
	Eval(*gql.Interface)
}
