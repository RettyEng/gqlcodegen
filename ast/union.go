package ast

import "github.com/RettyInc/gqlcodegen/gql"

type UnionInternalExpression interface {
	Eval(union *gql.Union)
}

type DefineUnionMemberExpression struct {
	TypeExp TypeRefExpression
}

func (d *DefineUnionMemberExpression) Eval(union *gql.Union) {
	union.Members = append(union.Members, d.TypeExp.Eval())
}
