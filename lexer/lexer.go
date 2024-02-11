package lexer

import "kat/token"

type Lexer struct {
	Input string
	Row   int
	Col   int
}

func New(input string) *Lexer {
	return &Lexer{Input: input, Row: 0, Col: 0}
}

func (l *Lexer) Next() token.Token {
	return token.Token{
		Row:   l.Row,
		Col:   l.Col,
		Value: string(l.Input[l.Col]),
		Type:  token.PLUS,
	}
}

func (l *Lexer) IsChar() bool {
	ch := l.Input[l.Col]
	return ch >= 'a' && ch <= 'z'
}
