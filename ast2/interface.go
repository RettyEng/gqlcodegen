package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type InterfaceInternalExpression interface {
	Eval(*gql.Interface)
}
type DefineInterfaceFieldExpression struct {
	NameExp        NameExpression
	TypeExp        TypeRefExpression
	DescriptionExp DescriptionExpression
	ArgsExp        []InputValueExpression
	DirectivesExp  []DirectiveExpression
}

func (d *DefineInterfaceFieldExpression) Eval(i *gql.Interface) {
	f := &gql.ObjectField{
		Name:        d.NameExp.Eval(),
		Type:        d.TypeExp.Eval(),
		Description: d.DescriptionExp.Eval(),
		Args:        evalInputValues(d.ArgsExp),
		Directives:  evalDirectives(d.DirectivesExp),
	}
	i.Fields = append(i.Fields, f)
}
