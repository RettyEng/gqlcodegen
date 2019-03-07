package gql

import (
	"strings"

	"github.com/RettyEng/gqlcodegen/ast/directive"
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

func NewTypeSystem() *TypeSystem {
	return &TypeSystem{
		Schema:           &Schema{},
		ScalarTypes:      map[string]*Scalar{},
		ObjectTypes:      map[string]*Object{},
		InterfaceTypes:   map[string]*Interface{},
		UnionTypes:       map[string]*Union{},
		EnumTypes:        map[string]*Enum{},
		InputObjectTypes: map[string]*InputObject{},
		Directives:       map[string]*Directive{},
	}
}

type Commentable interface {
	GetDescription() string
	GetDirectives() []*DirectiveRef
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

func (s *Scalar) GetDescription() string {
	return s.Description
}
func (s *Scalar) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type Object struct {
	Description string
	Name        string
	Implements  []*TypeRef
	Directives  []*DirectiveRef
	Fields      []*ObjectField
}

func (s *Object) GetDescription() string {
	return s.Description
}
func (s *Object) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type Interface struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Fields      []*ObjectField
}

func (s *Interface) GetDescription() string {
	return s.Description
}
func (s *Interface) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type Union struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Members     []*TypeRef
}

func (s *Union) GetDescription() string {
	return s.Description
}
func (s *Union) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type Enum struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Values      []*EnumValue
}

func (s *Enum) GetDescription() string {
	return s.Description
}
func (s *Enum) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type Directive struct {
	Description string
	Name        string
	Arguments   []*InputValue
	Location    []directive.Location
}

type EnumValue struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
}

func (s *EnumValue) GetDescription() string {
	return s.Description
}
func (s *EnumValue) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type ObjectField struct {
	Name        string
	Type        *TypeRef
	Description string
	Directives  []*DirectiveRef
	Args        []*InputValue
}

func (s *ObjectField) GetDescription() string {
	return s.Description
}
func (s *ObjectField) GetDirectives() []*DirectiveRef {
	return s.Directives
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
	Value() string
}
type ValueImpl struct {
	Val string
}

func (v *ValueImpl) Value() string {
	return v.Val
}

type List struct {
	ValueString string
	Child       []Value
}

func (l *List) Value() string {
	if l.ValueString == "null" {
		return "null"
	}
	var child []string
	for _, c := range l.Child {
		child = append(child, c.Value())
	}
	return "[" + strings.Join(child, ", ") + "]"
}

type InputObject struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	InputValue  []*InputValue
}

func (s *InputObject) GetDescription() string {
	return s.Description
}
func (s *InputObject) GetDirectives() []*DirectiveRef {
	return s.Directives
}

type InputValue struct {
	Description string
	Name        string
	Directives  []*DirectiveRef
	Type        *TypeRef
	Default     Value
}

func (s *InputValue) GetDescription() string {
	return s.Description
}
func (s *InputValue) GetDirectives() []*DirectiveRef {
	return s.Directives
}
