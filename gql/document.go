package gql

import (
	"strings"

	"github.com/RettyInc/gqlcodegen/ast2/directive"
)

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
	Name string
}

type ObjectField struct {
	Name        string
	Type        *TypeRef
	Description string
	Directives  []*DirectiveRef
	Args        []*InputValue
}

type DirectiveRef struct {
	Name string
	Args map[string]Value
}

type TypeRef struct {
	InnerType  *TypeRef
	Name       string
	IsNullable bool
}

type Value interface {
	Type() *TypeRef
	Value() string
	Children() []Value
}
type ValueImpl struct {
	typeRef *TypeRef
	value   string
}

func (v *ValueImpl) Type() *TypeRef {
	return v.typeRef
}
func (v *ValueImpl) Value() string {
	return v.Value()
}
func (v *ValueImpl) Children() []Value {
	return nil
}

type List struct {
	typeRef *TypeRef
	value   string
	child   []Value
}

func (l *List) Type() *TypeRef {
	return l.typeRef
}
func (l *List) Value() string {
	if l.value == "null" {
		return "null"
	}
	var child []string
	for _, c := range l.child {
		child = append(child, c.Value())
	}
	return "[" + strings.Join(child, ", ") + "]"
}
func (l *List) Children() []Value {
	return l.child
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
	Default     Value
}
