package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DirectiveExpression interface {
	Eval() *gql.DirectiveRef
}

type TypeRefExpression interface {
	Eval() *gql.TypeRef
}

type DescriptionExpression interface {
	Eval() string
}

type ListExpression interface {
	Eval() *gql.List
}

type ValueExpression interface {
	Eval() *gql.Value
}

type NameExpression interface {
	Eval() string
}