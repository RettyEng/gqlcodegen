package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type EnumInternalExpression interface {
	Eval(enum *gql.Enum)
}

type DefineEnumValueExpression struct {
	Description DescriptionExpression
	Name        NameExpression
	Directives  []DirectiveExpression
}

func (d *DefineEnumValueExpression) Eval(enum *gql.Enum) {
	enum.Values = append(enum.Values, &gql.EnumValue{
		Name:        d.Name.Eval(),
		Description: d.Description.Eval(),
		Directives:  evalDirectives(d.Directives),
	})
}
