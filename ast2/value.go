package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type TypeRefExpression interface {
	Eval() *gql.TypeRef
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

type InputValueExpression interface {
	Eval() *gql.InputValue
}
type InputValueExpressionImpl struct {
	descriptionExp DescriptionExpression
	nameExp        NameExpression
	typeExp        TypeRefExpression
	defaultValue   ValueExpression
	directives     []DirectiveExpression
}

func (exp *InputValueExpressionImpl) Eval(direc *gql.Directive) *gql.InputValue {
	return &gql.InputValue{
		Description: exp.descriptionExp.Eval(),
		Name:        exp.nameExp.Eval(),
		Type:        exp.typeExp.Eval(),
		Default:     exp.defaultValue.Eval(),
		Directives:  evalDirectives(exp.directives),
	}
}
