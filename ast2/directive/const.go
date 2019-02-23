package directive

type Location int

const (
	// Executable
	Query Location = iota
	Mutation
	Subscription
	Field
	FragmentDefinition
	FragmentSpread
	InlineFragment
	VariableDefinition

	// Type System
	Schema
	Scalar
	Object
	FieldDefinition
	ArgumentDefinition
	Interface
	Union
	Enum
	EnumValue
	InputObject
	InputFieldDefinition
)

//go:generate stringer -type=Location
