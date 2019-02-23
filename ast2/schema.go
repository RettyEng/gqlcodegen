package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type SchemaInternalExpression interface {
	Eval(schema *gql.Schema)
}

type DefineQueryExpression struct {
	Type TypeRefExpression
}

func (d *DefineQueryExpression) Eval(schema *gql.Schema) {
	schema.Query = d.Type.Eval()
}

type DefineMutationExpression struct {
	Type TypeRefExpression
}

func (d *DefineMutationExpression) Eval(schema *gql.Schema) {
	schema.Mutation = d.Type.Eval()
}

type DefineSubscriptionExpression struct {
	Type TypeRefExpression
}

func (d *DefineSubscriptionExpression) Eval(schema *gql.Schema) {
	schema.Subscription = d.Type.Eval()
}
