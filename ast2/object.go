package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type ObjectInternalExpression interface {
	Eval(object *gql.Object)
}

type ImplementExpression struct {
	TypeExp TypeRefExpression
}

func (e *ImplementExpression) Eval(object *gql.Object) {
	object.Implements = append(object.Implements, e.TypeExp.Eval())
}

type DefineFieldExpression struct {
	Name        NameExpression
	TypeRef     TypeRefExpression
	Description DescriptionExpression
	Directives  []DirectiveExpression
	Args        []InputValueExpression
}

func (e *DefineFieldExpression) Eval(object *gql.Object) {
	f := &gql.ObjectField{
		Name:        e.Name.Eval(),
		Type:        e.TypeRef.Eval(),
		Description: e.Description.Eval(),
		Directives:  evalDirectives(e.Directives),
		Args:        evalInputValues(e.Args),
	}
	object.Fields = append(object.Fields, f)
}
