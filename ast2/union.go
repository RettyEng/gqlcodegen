package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineUnionMemberExpression interface {
	Eval(union *gql.Union)
}

