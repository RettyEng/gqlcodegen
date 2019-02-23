package token

type Type int

const (
	TypeUnicodeBom Type = iota
	TypeWhiteSpace
	TypeLineTerminator
	TypeComma
	TypePunctuator
	TypeComment
	TypeName
	TypeIntVal
	TypeFloatVal
	TypeStrVal
)

//go:generate stringer -type=Type
