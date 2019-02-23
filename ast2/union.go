package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type UnionInternalExpression interface {
	Eval(union *gql.Union)
}

type DefineUnionMemberExpression struct {
	typeExp TypeRefExpression
}

func (d *DefineUnionMemberExpression) Eval(union *gql.Union) {
	union.Members = append(union.Members, d.typeExp.Eval())
}
