package lexer

type Token struct {
	Type  TokenType
	Value string
}

func schemaT() *Token {
	return &Token{Type: Schema}
}

func scalarT() *Token {
	return &Token{Type: Scalar}
}

func typeT() *Token {
	return &Token{Type: Type}
}

func lParenT() *Token {
	return &Token{Type: LParen}
}

func rParenT() *Token {
	return &Token{Type: RParen}
}

func lBraceT() *Token {
	return &Token{Type: LBrace}
}

func rBraceT() *Token {
	return &Token{Type: RBrace}
}

func lBracketT() *Token {
	return &Token{Type: LBracket}
}

func rBracketT() *Token {
	return &Token{Type: RBracket}

}

func newLineT() *Token {
	return &Token{Type: NewLine}
}

func colonT() *Token {
	return &Token{Type: Colon}
}

func commaT() *Token {
	return &Token{Type: Comma}
}

func notNullT() *Token {
	return &Token{Type: NotNull}
}

func eqT() *Token {
	return &Token{Type: Eq}
}

func numberT(value string) *Token {
	return &Token{Type: Number, Value: value}
}

func stringT(value string) *Token {
	return &Token{Type: String, Value: value}
}

func commentT(value string) *Token {
	return &Token{Type: Comment, Value: value}
}

func idT(value string) *Token {
	return &Token{Type: Id, Value: value}
}
