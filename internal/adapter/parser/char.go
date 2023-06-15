package parser

type Char rune

func (c Char) IsWhitespace() bool {
	return c == Char(' ')
}

func (c Char) IsQuote() bool {
	return c == Char('"')
}
