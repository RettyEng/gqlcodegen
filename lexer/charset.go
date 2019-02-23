package lexer

type _range struct {
	from rune
	to rune
}

type F struct {
	r rune
}

type Charset struct {
	_range []*_range
}

func Just(r rune) *Charset {
	return &Charset{
		_range: []*_range {
			{
				from: r,
				to: r,
			},
		},
	}
}

func From(r rune) *F {
	return &F{r}
}

func (f *F) To(r rune) *Charset {
	return &Charset{
		_range:[]*_range{
			{
				from: f.r,
				to: r,
			},
		},
	}
}

func Union(css ...*Charset) *Charset {
	var rs []*_range
	for _, cs := range css {
		rs = append(rs, cs._range...)
	}
	return &Charset{ rs }
}

func (c *Charset) Contains(r rune) bool {
	for _, _range := range c._range {
		if r < _range.from ||_range.to < r {
			return false
		}
	}
	return true
}
