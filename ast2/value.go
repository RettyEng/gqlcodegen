package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type TypeRefExpression interface {
	Eval() *gql.TypeRef
}
type TypeRefExpressionImpl struct {
	innerType  TypeRefExpression
	isNullable bool
	name       NameExpression
}

func (exp *TypeRefExpressionImpl) Eval() *gql.TypeRef {
	var inner *gql.TypeRef = nil
	if exp.innerType != nil {
		inner = exp.innerType.Eval()
	}
	return &gql.TypeRef{
		InnerType:  inner,
		Name:       exp.name.Eval(),
		IsNullable: exp.isNullable,
	}
}

type ValueExpression interface {
	Eval() gql.Value
}
type ValueExpressionImpl struct {
	value   string
	typeRef TypeRefExpression
}

func (exp *ValueExpressionImpl) Eval() gql.Value {
	return &gql.ValueImpl{exp.typeRef.Eval(), exp.value}
}

type ListValueExpressionImpl struct {
	typeRef  TypeRefExpression
	children []ValueExpression
}

func (exp *ListValueExpressionImpl) Eval() gql.Value {
	var child []gql.Value
	for _, e := range exp.children {
		child = append(child, e.Eval())
	}
	return &gql.List{}
}

type NameExpression interface {
	Eval() string
}
type NameExpressionImpl struct {
	Name string
}

func (exp *NameExpressionImpl) Eval() string {
	return exp.Name
}

type DescriptionExpression interface {
	Eval() string
}
type DescriptionExpressionImpl struct {
	description string
}

func (exp *DescriptionExpressionImpl) Eval() string {
	return exp.description
}

type DirectiveExpression interface {
	Eval() *gql.DirectiveRef
}
type DirectiveExpressionImpl struct {
	name string
	args map[string]ValueExpression
}

func (exp *DirectiveExpressionImpl) Eval() *gql.DirectiveRef {
	args := map[string]gql.Value{}
	for name, v := range exp.args {
		args[name] = v.Eval()
	}
	return &gql.DirectiveRef{
		Name: exp.name,
		Args: args,
	}
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
