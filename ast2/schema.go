package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type SchemaInternalExpression interface {
	Eval(schema *gql.Schema)
}

type DefineQueryExpression struct {
	typeExp TypeRefExpression
}

func (d *DefineQueryExpression) Eval(schema *gql.Schema) {
	schema.Query = d.typeExp.Eval()
}

type DefineMutationExpression struct {
	typeExp TypeRefExpression
}

func (d *DefineMutationExpression) Eval(schema *gql.Schema) {
	schema.Mutation = d.typeExp.Eval()
}

type DefineSubscriptionExpression struct {
	typeExp TypeRefExpression
}

func (d *DefineSubscriptionExpression) Eval(schema *gql.Schema) {
	schema.Subscription = d.typeExp.Eval()
}
