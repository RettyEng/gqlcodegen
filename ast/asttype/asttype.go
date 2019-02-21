package asttype

type AstType int

const (
	Root AstType = iota

	SchemaDef
	SchemaQueryDef

	ScalarDef

	TypeDef
	TypeFieldDef
	FieldArgDef
	FieldArgDefault

	EnumDef
	EnumEntryDef

	TypeRef
)

//go:generate stringer -type=AstType
