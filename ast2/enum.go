package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineEnumValuesExpression interface {
	Eval(enum *gql.Enum)
}
