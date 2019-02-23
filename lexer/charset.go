package lexer

type Matcher interface {
	HeadMatches(runes []rune) bool
	MatchCount(runes []rune) int
}

/*********************************************************
Not
 *********************************************************/

type not struct {
	matcher Matcher
}

func Not(m Matcher) Matcher {
	return &not{m}
}

func (n *not) HeadMatches(rs []rune) bool {
	if len(rs) == 0 {
		return false
	}
	return !n.matcher.HeadMatches(rs)
}

func (n *not) MatchCount(rs []rune) int {
	c := 0
	for {
		if n.matcher.HeadMatches(rs[c:]) {
			return c
		}
		c++
	}
}

/*********************************************************
Union
 *********************************************************/

type union struct {
	matchers []Matcher
}

func Union(ms ...Matcher) Matcher {
	return &union{ms}
}

func StrUnion(strs ...string) Matcher {
	ms := make([]Matcher, len(strs))
	for i, str := range strs {
		ms[i] = Str(str)
	}
	return Union(ms...)
}

func (u *union) HeadMatches(rs []rune) bool {
	for _, m := range u.matchers {
		if m.HeadMatches(rs) {
			return true
		}
	}
	return false
}

func (u *union) MatchCount(runes []rune) int {
	for _, m := range u.matchers {
		if m.HeadMatches(runes) {
			return m.MatchCount(runes)
		}
	}
	return 0
}

/*********************************************************
Str
 *********************************************************/

type strMatch struct {
	s string
}

func Str(str string) Matcher {
	return &strMatch{str}
}

func (sm *strMatch) HeadMatches(rs []rune) bool {
	if len([]rune(sm.s)) > len(rs) {
		return false
	}
	for i, r := range []rune(sm.s) {
		if rs[i] != r {
			return false
		}
	}
	return true
}

func (sm *strMatch) MatchCount(runes []rune) int {
	if sm.HeadMatches(runes) {
		return len([]rune(sm.s))
	}
	return 0
}

/*********************************************************
Charset
 *********************************************************/

type charset struct {
	from rune
	to   rune
}

type F struct {
	r rune
}

func From(r rune) *F {
	return &F{r}
}

func (f *F) To(r rune) Matcher {
	return &charset{
		from: f.r,
		to:   r,
	}
}

func (c *charset) HeadMatches(rs []rune) bool {
	if len(rs) == 0 {
		return false
	}
	r := rs[0]
	return c.from <= r && r <= c.to
}

func (c *charset) MatchCount(rs []rune) int {
	if c.HeadMatches(rs) {
		return 1
	}
	return 0
}

/*********************************************************
Charsets
 *********************************************************/

var (
	digit        = From('0').To('9')
	nonZeroDigit = From('1').To('9')

	unicodeBom = Str("\ufeff")
	whiteSpace = StrUnion("\u0009", "\u0020")

	lineTerminator = StrUnion("\n", "\r\n", "\r")

	comma = Str(",")

	punctuator = StrUnion(
		"!", "$", "(", ")", "...", ":", "=", "@", "[", "]", "{", "}", "|",
	)

	commentHead = Str("#")
	commentTail = Not(lineTerminator)

	nameHead = Union(Str("_"), From('A').To('Z'), From('a').To('z'))
	nameTail = Union(nameHead, digit)

	negative = Str("-")
	intVal   = digit

	fractionalPartHead = Str(".")
	exponent           = StrUnion("e", "E")
	sign               = Union(Str("+"), negative)

	strStart             = Str("\"")
	strEnd               = Str("\"")
	strUnicodeEscapeHead = Str(`\u`)
	hex                  = Union(From(0).To(9), From('A').To('F'), From('a').To('f'))
	strEscape            = StrUnion(
		`\\`, `\"`, `\/`, `\b`, `\f`, `\n`, `\r`, `\t`,
	)
	strChar = Union(
		strEscape,
		Not(
			Union(Str("\""), Str("\\"), lineTerminator),
		),
	)

	blockStrStart = Str(`"""`)
	blockStrEnd   = Str(`"""`)
	blockStrChar  = Union(Str(`\"""`), Not(Str(`"""`)))
)
