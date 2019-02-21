package lexer

import (
	"bufio"
)

type Lexer struct {
	buff *bufio.Reader
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{buff: r}
}

func (t*Lexer) takeRunes() []rune {
	runes := make([]rune, 0)
	for {
		var r rune
		var e error

		if r, _, e = t.buff.ReadRune(); e != nil {
			break
		}
		runes = append(runes, r)
	}
	return runes
}

func (t *Lexer) Tokenize() []string {
	tokens := make([]string, 0)
	runes := t.takeRunes()
	cur := 0

	for cur < len(runes) {
		var token []rune
		switch r := runes[cur]; r {
		case ' ', '\t':
			cur++
			continue
		case '(', ')', '{', '}', '[', ']', '\n', ':', '!', '=':
			token = []rune{r}
		case '"':
			token = takeString(runes[cur:])
		case '#':
			token = takeComment(runes[cur:])
		default:
			token = takeWhileSeparator(runes[cur:])
		}
		cur += len(token)
		tokens = append(tokens, string(token))
	}
	return tokens
}

func takeString(rs []rune) []rune {
	str := make([]rune, 0)
	cur := 0
	isEscaping := true
	r := rs[cur]
	for {
		str = append(str, r)
		if r == '"' && !isEscaping {
			break
		}
		if !isEscaping && r == '\\' {
			isEscaping = true
		} else {
			isEscaping = false
		}
		cur++
		r = rs[cur]
	}
	return str
}

func takeComment(rs []rune) []rune {
	comment := make([]rune, 0)
	for cur := 0; rs[cur] != '\n'; cur++ {
		comment = append(comment, rs[cur])
	}
	return comment
}

func takeWhileSeparator(rs []rune) []rune {
	token := make([]rune, 0)
	separators := map[rune]struct{} {
		' ': {},
		'\t': {},
		'(': {},
		')': {},
		'{': {},
		'}': {},
		'[': {},
		']': {},
		'\n': {},
		':': {},
		'!': {},
		'=': {},
	}

	cur := 0
	for {
		r := rs[cur]
		if _, contains := separators[r]; contains {
			break
		}
		token = append(token, r)
		cur++
	}
	return token
}