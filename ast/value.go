package ast

import "github.com/RettyInc/gqlcodegen/gql"

type TypeRefExpression interface {
	Eval() *gql.TypeRef
}
type TypeRefExpressionImpl struct {
	InnerType  TypeRefExpression
	IsNullable bool
	Name       NameExpression
}

func (exp *TypeRefExpressionImpl) Eval() *gql.TypeRef {
	var inner *gql.TypeRef = nil
	if exp.InnerType != nil {
		inner = exp.InnerType.Eval()
	}
	return &gql.TypeRef{
		InnerType:  inner,
		Name:       exp.Name.Eval(),
		IsNullable: exp.IsNullable,
	}
}

type ValueExpression interface {
	Eval() gql.Value
}
type ValueExpressionImpl struct {
	Value string
}

func (exp *ValueExpressionImpl) Eval() gql.Value {
	return &gql.ValueImpl{
		Val: exp.Value,
	}
}

type ListValueExpressionImpl struct {
	Children []ValueExpression
}

func (exp *ListValueExpressionImpl) Eval() gql.Value {
	var child []gql.Value
	for _, e := range exp.Children {
		child = append(child, e.Eval())
	}
	return &gql.List{
		"[]",
		child,
	}
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
	Description string
}

func (exp *DescriptionExpressionImpl) Eval() string {
	return exp.Description
}

type EmptyDescription struct{}

func (exp *EmptyDescription) Eval() string {
	return ""
}

type DirectiveExpression interface {
	Eval() *gql.DirectiveRef
}
type DirectiveExpressionImpl struct {
	Name string
	Args map[string]ValueExpression
}

func (exp *DirectiveExpressionImpl) Eval() *gql.DirectiveRef {
	args := map[string]gql.Value{}
	for name, v := range exp.Args {
		args[name] = v.Eval()
	}
	return &gql.DirectiveRef{
		Name: exp.Name,
		Args: args,
	}
}

type InputValueExpression interface {
	Eval() *gql.InputValue
}
type InputValueExpressionImpl struct {
	Description  DescriptionExpression
	Name         NameExpression
	Type         TypeRefExpression
	DefaultValue ValueExpression
	Directives   []DirectiveExpression
}

func (exp *InputValueExpressionImpl) Eval() *gql.InputValue {
	var value gql.Value = nil
	if exp.DefaultValue != nil {
		value = exp.DefaultValue.Eval()
	}
	return &gql.InputValue{
		Description: exp.Description.Eval(),
		Name:        exp.Name.Eval(),
		Type:        exp.Type.Eval(),
		Default:     value,
		Directives:  evalDirectives(exp.Directives),
	}
}
