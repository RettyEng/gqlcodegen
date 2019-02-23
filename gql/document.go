package gql

import "github.com/RettyInc/gqlcodegen/ast2/directive"

type TypeSystem struct {
	Schema           *Schema
	ScalarTypes      map[string]*Scalar
	ObjectTypes      map[string]*Object
	InterfaceTypes   map[string]*Interface
	UnionTypes       map[string]*Union
	EnumTypes        map[string]*Enum
	InputObjectTypes map[string]*InputObject
	Directives       map[string]*Directive
}

type Schema struct {
	Directives   []*DirectiveRef
	Query        *TypeRef
	Mutation     *TypeRef
	Subscription *TypeRef
}

type Scalar struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
}

type Object struct {
	Description string
	Name        string
	Implements  []*TypeRef
	Directives  []*DirectiveRef
	Fields      []*ObjectField
}

type Interface struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Fields      []*ObjectField
}

type Union struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Members     []*TypeRef
}

type Enum struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Values      []*EnumValue
}

type Directive struct {
	Description string
	Name        string
	Arguments   []*InputValue
	Location    []directive.Location
}

type EnumValue struct {
}

type ObjectField struct {
	Name        string
	Type        *TypeRef
	Description string
	Directives  []*DirectiveRef
	Args        []*InputValue
}

type DirectiveRef struct {
}

type TypeRef struct {
}

type List struct {
}

type Value struct {
}

type InputObject struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	InputValue  []*InputValue
}

type InputValue struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Type        *TypeRef
	Default     *Value
}

type Name string
