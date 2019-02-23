package lexer

import (
	"bufio"
	"io"
	"log"

	"github.com/RettyInc/gqlcodegen/lexer/token"
)

type Lexer struct {
	scanner *Scanner
}

func NewLexer(r io.Reader) *Lexer {
	buff := bufio.NewReader(r)
	var runes []rune
	for {
		r, _, e := buff.ReadRune()
		if e != nil {
			break
		}
		runes = append(runes, r)
	}
	return &Lexer{NewScanner(runes)}
}

func (l *Lexer) next() *token.Token {
	s := l.scanner
	takeWhileAndAppend := func(t token.Type, m ...Matcher) *token.Token {
		v, c, l := s.TakeWhileMatch(m[0])
		for _, m := range m[1:] {
			tail, _, _ := s.TakeWhileMatch(m)
			v += tail
		}
		return token.NewToken(t, v, c, l)
	}
	takeAndAppend := func(t token.Type, m ...Matcher) *token.Token {
		v, c, l := s.Take(m[0])
		return token.NewToken(t, v, c, l)
	}

	if !l.scanner.HasNext() {
		return nil
	}

	if s.StartsWith(unicodeBom) {
		return takeAndAppend(token.TypeUnicodeBom, unicodeBom)
	}

	if s.StartsWith(whiteSpace) {
		return takeAndAppend(token.TypeWhiteSpace, whiteSpace)
	}

	if s.StartsWith(lineTerminator) {
		return takeAndAppend(token.TypeLineTerminator, lineTerminator)
	}

	if s.StartsWith(comma) {
		return takeAndAppend(token.TypeComma, comma)
	}

	if s.StartsWith(punctuator) {
		return takeAndAppend(token.TypePunctuator, punctuator)
	}

	if s.StartsWith(commentHead) {
		return takeWhileAndAppend(token.TypeComment, commentHead, commentTail)
	}

	if s.StartsWith(nameHead) {
		return takeWhileAndAppend(token.TypeName, nameHead, nameTail)
	}

	if s.StartsWith(negative) || s.StartsWith(intVal) {
		return l.takeNumber()
	}

	if s.StartsWith(blockStrStart) {
		return l.takeBlockString()
	}

	if s.StartsWith(strStart) {
		return l.takeString()
	}

	log.Fatalf(
		"unexpected token %c at line %d, col %d",
		s.runes[0], s.line, s.col,
	)
	return nil
}

func (l *Lexer) Next() *token.Token {
	ignored := map[token.Type]struct{}{
		token.TypeUnicodeBom:     {},
		token.TypeWhiteSpace:     {},
		token.TypeLineTerminator: {},
		token.TypeComment:        {},
		token.TypeComma:          {},
	}

	for {
		t := l.next()
		if t == nil {
			return nil
		}
		if _, isIgnored := ignored[t.Type()]; !isIgnored {
			return t
		}
	}
}

func (l *Lexer) takeNumber() *token.Token {
	s := l.scanner
	negativeSign, line, col := s.TakeWhileMatch(negative)
	fatal := func() {
		log.Fatalf("illegal int token at line %d, col %d", line, col)
	}
	intPart, _, _ := s.TakeWhileMatch(intVal)
	if len([]rune(intPart)) == 0 {
		fatal()
	}
	if rs := []rune(intPart); rs[0] == '0' && len(rs) > 1 {
		fatal()
	}
	fracHead, _, _ := s.TakeWhileMatch(fractionalPartHead)
	fracPart := ""
	if len([]rune(fracHead)) != 0 {
		fracPart, _, _ = s.TakeWhileMatch(intVal)
	}
	exponentHead, _, _ := s.TakeWhileMatch(exponent)
	exponentPart := ""
	if len([]rune(exponentHead)) != 0 {
		exponentPart, _, _ = s.TakeWhileMatch(sign)
		n, _, _ := s.TakeWhileMatch(intVal)
		exponentPart += n
	}

	if len([]rune(fracHead)) != 0 || len([]rune(exponentHead)) != 0 {
		return token.NewToken(
			token.TypeFloatVal,
			negativeSign+intPart+fracHead+fracPart+exponentHead+exponentPart,
			line, col,
		)
	}
	return token.NewToken(token.TypeIntVal, negativeSign+intPart, line, col)
}

func (l *Lexer) takeBlockString() *token.Token {
	s := l.scanner
	value, line, col := s.Take(blockStrStart)
	v, _, _ := s.Take(blockStrChar)
	value += v
	if !s.StartsWith(blockStrEnd) {
		log.Fatalf("illegal string at line %d, col %d", line, col)
	}
	v, _, _ = s.Take(blockStrEnd)
	value += v
	return token.NewToken(token.TypeStrVal, value, line, col)
}

func (l *Lexer) takeString() *token.Token {
	s := l.scanner
	value, line, col := s.TakeWhileMatch(strStart)
	for {
		v, _, _ := s.TakeWhileMatch(strChar)
		value += v
		if s.StartsWith(strUnicodeEscapeHead) {
			u, ul, uc := s.Take(strUnicodeEscapeHead)
			value += u
			for i := 0; i < 4; i++ {
				if !s.StartsWith(hex) {
					log.Fatalf("illegal unicode escape at line %d, col %d", ul, uc)
				}
				h, _, _ := s.Take(hex)
				value += h
			}
			continue
		}
		if s.StartsWith(strEnd) {
			v, _, _ = s.Take(strEnd)
			value += v
			break
		}
		log.Fatalf("illegal string at line %d, col %d", line, col)
	}
	return token.NewToken(token.TypeStrVal, value, line, col)
}
