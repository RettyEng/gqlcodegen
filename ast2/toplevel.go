package ast2

import "github.com/RettyInc/gqlcodegen/gql"


type TopLevel struct {
	expressions []TopLevelExpression
}
type TopLevelExpression interface {
	Eval(system *gql.TypeSystem)
}


type DefineSchemaExpression struct {
	directiveExpressions []DirectiveExpression
	expressions []DefineSchemaFieldExpression
}

type ExtendSchemaExpression struct {
	directiveExpressions []DirectiveExpression
	expressions []DefineSchemaFieldExpression
}

type DefineScalarExpression struct {
	descriptionExpression DescriptionExpression
	nameExpression        NameExpression
	directiveExpressions  []DirectiveExpression
}

type ExtendScalarExpression struct {
	nameExpression        NameExpression
	directiveExpressions  []DirectiveExpression
}

type DefineObjectExpression struct {
	descriptionExpression  DescriptionExpression
	nameExpression         NameExpression
	implementsExpressions  []ImplementsExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineFieldExpression
}

type ExtendObjectExpression struct {
	nameExpression         NameExpression
	implementsExpressions  []ImplementsExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineFieldExpression
}

type DefineInterfaceExpression struct {
	descriptionExpression  DescriptionExpression
	nameExpression         NameExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineInterfaceFieldExpression
}

type ExtendInterfaceExpression struct {
	nameExpression         NameExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineInterfaceFieldExpression
}

type DefineUnionExpression struct {
	descriptionExpression       DescriptionExpression
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineUnionMemberExpression []DefineUnionMemberExpression
}

type ExtendUnionExpression struct {
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineUnionMemberExpression []DefineUnionMemberExpression
}

type DefineEnumExpression struct {
	descriptionExpression       DescriptionExpression
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineEnumValuesExpressions []DefineEnumValuesExpression
}

type ExtendEnumExpression struct {
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineEnumValuesExpressions []DefineEnumValuesExpression
}

type DefineInputObjectExpression struct {
	descriptionExpression             DescriptionExpression
	nameExpression                    NameExpression
	directiveExpressions              []DirectiveExpression
	defineInputObjectFieldExpressions []DefineInputObjectFieldExpression
}

type ExtendInputObjectExpression struct {
	nameExpression                    NameExpression
	directiveExpressions              []DirectiveExpression
	defineInputObjectFieldExpressions []DefineInputObjectFieldExpression
}

type DirectiveDefinition struct {
	descriptionExpression             DescriptionExpression
	nameExpression                    NameExpression
	defineArgsExpression []DefineDirectiveArgsExpression
}
