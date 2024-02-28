package lexer

import (
	"kat/token"
)

type Lexer struct {
	Col    int
	Line   int
	Offset int
	Input  []byte
}

func New(input []byte) *Lexer {
	return &Lexer{Col: -1, Line: 0, Offset: 0, Input: input}
}

func (l *Lexer) MakeToken(col int, val string, tokenType token.TokenType) token.Token {
	return token.Token{
		Row:   l.Line,
		Col:   col - l.Offset,
		Value: val,
		Type:  tokenType,
	}
}

func (l *Lexer) NextToken() token.Token {
	l.NextChar()
	l.SkipWhitespace()

	ch := l.Char()
	var t token.Token

	switch ch {

	case '+':
		if l.PeekChar() == '+' {
			col := l.Col
			l.NextChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.PLUSPLUS)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.PLUS)
		}

	case '-':
		if l.PeekChar() == '-' {
			col := l.Col
			l.NextChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.MINUSMINUS)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.MINUS)
		}

	case '=':
		if l.PeekChar() == '=' {
			col := l.Col
			l.NextChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.EQUALEQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.EQUAL)
		}

	case '<':
		if l.PeekChar() == '=' {
			col := l.Col
			l.NextChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.LESSEQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.LESS)
		}

	case '>':
		if l.PeekChar() == '=' {
			col := l.Col
			l.NextChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.GREATEREQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.GREATER)
		}

	case '*':
		t = l.MakeToken(l.Col, string(ch), token.MULTIPLY)

	case '/':
		t = l.MakeToken(l.Col, string(ch), token.DIVIDE)

	case '%':
		t = l.MakeToken(l.Col, string(ch), token.MODULO)

	case '[':
		t = l.MakeToken(l.Col, string(ch), token.LBRACKET)

	case ']':
		t = l.MakeToken(l.Col, string(ch), token.RBRACKET)

	case '{':
		t = l.MakeToken(l.Col, string(ch), token.LBRACE)

	case '}':
		t = l.MakeToken(l.Col, string(ch), token.RBRACE)

	case '(':
		t = l.MakeToken(l.Col, string(ch), token.LPAREN)

	case ')':
		t = l.MakeToken(l.Col, string(ch), token.RPAREN)

	case '!':
		t = l.MakeToken(l.Col, string(ch), token.BANG)

	case '?':
		t = l.MakeToken(l.Col, string(ch), token.QUESTION)

	case '"':
		col := l.Col
		str := l.MakeString()
		t = l.MakeToken(col, string(str), token.STRING)

	case ':':
		t = l.MakeToken(l.Col, string(ch), token.COLON)

	case ',':
		t = l.MakeToken(l.Col, string(ch), token.COMMA)

	case ';':
		t = l.MakeToken(l.Col, string(ch), token.SEMICOLON)

	case '.':
		t = l.MakeToken(l.Col, string(ch), token.DOT)

	case '\n':
		t = l.MakeToken(l.Col, "\\n", token.EOL)
		l.Offset = l.Col + 1
		l.Line++

	case 0:
		t = l.MakeToken(l.Col, "EOF", token.EOF)

	default:
		if l.IsDigit(ch) {
			col := l.Col
			dig := l.MakeDigit()
			t = l.MakeToken(col, string(dig), token.DIGIT)
		} else if l.IsChar(ch) {
			col := l.Col
			unknown := l.MakeIdentifier()
			symbol := string(unknown)
			t = l.MakeToken(col, symbol, token.Symbol(symbol))
		} else {
			col := l.Col
			invalid := l.MakeInvalid()
			t = l.MakeToken(col, string(invalid), token.INVALID)
		}
	}

	return t
}

func (l *Lexer) Char() byte {
	if l.Col < len(l.Input) {
		return l.Input[l.Col]
	}

	return 0
}

func (l *Lexer) NextChar() {
	l.Col++
}

func (l *Lexer) PeekChar() byte {
	if l.Col+1 < len(l.Input) {
		return l.Input[l.Col+1]
	}

	return 0
}

func (l *Lexer) IsChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) IsDigit(ch byte) bool {
	return ch >= '0' && ch <= '9' || ch == '.'
}

func (l *Lexer) IsAlphaNum(ch byte) bool {
	return l.IsChar(ch) || l.IsDigit(ch) || ch == '_'
}

func (l *Lexer) IsWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func (l *Lexer) IsEndOfString() bool {
	return l.Char() == '"' || l.Char() == 0
}

func (l *Lexer) SkipWhitespace() {
	for l.Col < len(l.Input) && l.IsWhitespace(l.Input[l.Col]) {
		l.NextChar()
	}
}

func (l *Lexer) MakeString() []byte {
	start := l.Col
	l.NextChar()

	for !l.IsEndOfString() {
		l.NextChar()
	}

	end := l.Col
	return l.Input[start : end+1]
}

func (l *Lexer) MakeDigit() []byte {
	start := l.Col
	l.NextChar()

	for l.IsDigit(l.Char()) {
		l.NextChar()
	}

	end := l.Col
	l.Col--

	return l.Input[start:end]
}

func (l *Lexer) MakeIdentifier() []byte {
	start := l.Col
	l.NextChar()

	for l.IsAlphaNum(l.Char()) {
		l.NextChar()
	}

	end := l.Col
	l.Col--
	return l.Input[start:end]
}

func (l *Lexer) MakeInvalid() []byte {
	start := l.Col

	for l.Char() != 0 && !l.IsWhitespace(l.Char()) {
		l.NextChar()
	}

	end := l.Col
	l.Col--
	return l.Input[start:end]
}
