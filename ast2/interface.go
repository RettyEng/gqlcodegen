package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type InterfaceInternalExpression interface {
	Eval(*gql.Interface)
}
type DefineInterfaceFieldExpression struct {
	nameExp        NameExpression
	typeExp        TypeRefExpression
	descriptionExp DescriptionExpression
	argsExp        []InputValueExpression
	directivesExp  []DirectiveExpression
}

func (d *DefineInterfaceFieldExpression) Eval(i *gql.Interface) {
	f := &gql.ObjectField{
		Name:        d.nameExp.Eval(),
		Type:        d.typeExp.Eval(),
		Description: d.descriptionExp.Eval(),
		Args:        evalInputValues(d.argsExp),
		Directives:  evalDirectives(d.directivesExp),
	}
	i.Fields = append(i.Fields, f)
}
