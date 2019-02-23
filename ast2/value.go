package ast2

import "github.com/RettyInc/gqlcodegen/gql"


type TypeRefExpression interface {
	Eval() *gql.TypeRef
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

type DescriptionExpression interface {
	Eval() string
}

type DirectiveExpression interface {
	Eval() *gql.DirectiveRef
}
