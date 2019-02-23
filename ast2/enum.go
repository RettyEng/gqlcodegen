package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type EnumInternalExpression interface {
	Eval(enum *gql.Enum)
}

type DefineEnumValueExpression struct {
	name NameExpression
}

func (d *DefineEnumValueExpression) Eval(enum *gql.Enum) {
	enum.Name = d.name.Eval()
}
