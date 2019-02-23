package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type ObjectInternalExpression interface {
	Eval(object *gql.Object)
}

type ImplementExpression struct {
	typeExp TypeRefExpression
}

func (e *ImplementExpression) Eval(object *gql.Object) {
	object.Implements = append(object.Implements, e.typeExp.Eval())
}

type DefineFieldExpression struct {
	name        NameExpression
	typeRef     TypeRefExpression
	description DescriptionExpression
	directives  []DirectiveExpression
	args        []InputValueExpression
}

func (e *DefineFieldExpression) Eval(object *gql.Object) {
	f := &gql.ObjectField{
		Name:        e.name.Eval(),
		Type:        e.typeRef.Eval(),
		Description: e.description.Eval(),
		Directives:  evalDirectives(e.directives),
		Args:        evalInputValues(e.args),
	}
	object.Fields = append(object.Fields, f)
}
