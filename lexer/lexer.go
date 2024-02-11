package lexer

import (
	"kat/token"
)

type Lexer struct {
	Col    int
	Line   int
	Input  []byte
	Offset int
}

func New(input []byte) *Lexer {
	return &Lexer{Input: input, Line: 0, Col: 0}
}

func (l *Lexer) MakeToken(col int, val string, tokenType token.TokenType) token.Token {
	return token.Token{
		Row:   l.Line,
		Col:   col - l.Offset,
		Value: val,
		Type:  tokenType,
	}
}

func (l *Lexer) Token() token.Token {
	l.SkipWhitespace()

	ch := l.Char()
	var t token.Token

	switch ch {

	case '+':
		if l.PeekChar() == '+' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.PLUSPLUS)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.PLUS)
		}

	case '-':
		if l.PeekChar() == '-' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.MINUSMINUS)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.MINUS)
		}

	case '=':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.EQUALEQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.EQUAL)
		}

	case '<':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.LESSEQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.LESS)
		}

	case '>':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.GREATEREQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.GREATER)
		}

	case '*':
		t = l.MakeToken(l.Col, string(ch), token.MULTIPLY)

	case '/':
		t = l.MakeToken(l.Col, string(ch), token.DIVIDE)

	case '%':
		t = l.MakeToken(l.Col, string(ch), token.MODULE)

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
			unknown := l.MakeUnknown()
			symbol := string(unknown)
			t = l.MakeToken(col, symbol, token.Symbol(symbol))
		} else {
			col := l.Col
			unknown := l.MakeUnknown()
			t = l.MakeToken(col, string(unknown), token.UNKNOWN)
		}
	}

	l.AdvanceChar()
	return t
}

func (l *Lexer) Char() byte {
	if l.Col < len(l.Input) {
		return l.Input[l.Col]
	}

	return 0
}

func (l *Lexer) AdvanceChar() {
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
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) IsAlphaNum(ch byte) bool {
	return l.IsChar(ch) || l.IsDigit(ch) || ch == '_'
}

func (l *Lexer) IsWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

func (l *Lexer) IsQuote() bool {
	return l.Char() == '"'
}

func (l *Lexer) SkipWhitespace() {
	for l.Col < len(l.Input) && l.IsWhitespace(l.Input[l.Col]) {
		l.AdvanceChar()
	}
}

func (l *Lexer) MakeString() []byte {
	l.AdvanceChar()
	start := l.Col

	for !l.IsQuote() {
		l.AdvanceChar()
	}

	end := l.Col

	return l.Input[start:end]
}

func (l *Lexer) MakeDigit() []byte {
	start := l.Col
	l.AdvanceChar()

	for l.IsDigit(l.Char()) {
		l.AdvanceChar()
	}

	end := l.Col
	l.Col--

	return l.Input[start:end]
}

func (l *Lexer) MakeUnknown() []byte {
	start := l.Col
	l.AdvanceChar()

	for l.IsAlphaNum(l.Char()) {
		l.AdvanceChar()
	}

	end := l.Col
	l.Col--
	return l.Input[start:end]
}
