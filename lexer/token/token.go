package token

type Token struct {
	tokenType Type
	value string
	line int
	col int
}

func NewToken(t Type, v string, l, c int) *Token {
	return &Token{t, v, l, c}
}

func (t *Token) Type() Type {
	return t.tokenType
}

func (t *Token) Value() string {
	return t.value
}

func (t *Token) LineCol() (int, int) {
	return t.line, t.col
}
