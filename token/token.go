package token

import (
	"fmt"
)

var Precedence = struct {
	ASSIGNMENT  int
	CONDITIONAL int
	COMPARISON  int
	SUM         int
	PRODUCT     int
	EXPONENT    int
	PREFIX      int
	POSTFIX     int
	CALL        int
}{
	ASSIGNMENT:  0,
	CONDITIONAL: 1,
	COMPARISON:  2,
	SUM:         3,
	PRODUCT:     4,
	EXPONENT:    5,
	PREFIX:      6,
	POSTFIX:     7,
	CALL:        8,
}

type TokenType string

type Token struct {
	Row   int
	Col   int
	Value string
	Type  TokenType
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token{ Line: %d, Col: %d, Type: %s, Value: `%s` }",
		t.Row, t.Col, t.Type, t.Value,
	)
}

const (
	// Single character
	PLUS      = "PLUS"      // +
	MINUS     = "MINUS"     // -
	NEGATE    = "NEGATE"    // -
	BANG      = "BANG"      // !
	QUESTION  = "QUESTION"  // ?
	MULTIPLY  = "MULTIPLY " // *
	DIVIDE    = "DIVIDE"    // /
	MODULO    = "MODULO"    // %
	EQUAL     = "EQUAL"     // =
	LESS      = "LESS"      // <
	GREATER   = "GREATER"   // >
	LBRACKET  = "LBRACKET"  // [
	RBRACKET  = "RBRACKET"  // ]
	LBRACE    = "LBRACE"    // {
	RBRACE    = "RBRACE"    // }
	COLON     = "COLON"     // :
	LPAREN    = "LPAREN"    // (
	RPAREN    = "RPAREN"    // )
	COMMA     = "COMMA"     // ,
	SEMICOLON = "SEMICOLON" // ;
	DOT       = "DOT"       // .

	// Double character
	PLUSPLUS     = "PLUSPLUS"     // ++
	MINUSMINUS   = "MINUSMINUS"   // --
	EQUALEQUAL   = "EQUALEQUAL"   // ==
	GREATEREQUAL = "GREATEREQUAL" // >=
	LESSEQUAL    = "LESSEQUAL"    // <=

	// Literal
	STRING = "STRING"
	DIGIT  = "DIGIT"

	// Keyword

	TRUE       = "TRUE"       // true
	FALSE      = "FALSE"      // false
	LET        = "LET"        // let
	CONST      = "CONST"      // const
	IF         = "IF"         // if
	ELSE       = "ELSE"       // else
	FOR        = "FOR"        // for
	SELF       = "SELF"       // self
	IMPORT     = "IMPORT"     // import
	STRUCT     = "STRUCT"     // struct
	FUNCTION   = "FUNCTION"   // fn
	IDENTIFIER = "IDENTIFIER" // any

	// Special
	EOL     = "EOL"     // End of line
	EOF     = "EOF"     // End of file
	INVALID = "INVALID" // End of file
)

func Symbol(key string) TokenType {
	keywords := map[string]TokenType{
		"let":    LET,
		"const":  CONST,
		"if":     IF,
		"else":   ELSE,
		"for":    FOR,
		"self":   SELF,
		"import": IMPORT,
		"struct": STRUCT,
		"fn":     FUNCTION,
		"true":   TRUE,
		"false":  FALSE,
	}

	keyword, ok := keywords[key]

	if ok {
		return keyword
	}

	return IDENTIFIER
}

func GetPrecedence(tok Token) int {
	precedences := map[TokenType]int{
		EQUAL:        Precedence.ASSIGNMENT,
		LESS:         Precedence.COMPARISON,
		GREATER:      Precedence.COMPARISON,
		LESSEQUAL:    Precedence.COMPARISON,
		GREATEREQUAL: Precedence.COMPARISON,
		PLUS:         Precedence.SUM,
		MINUS:        Precedence.SUM,
		MULTIPLY:     Precedence.PRODUCT,
		DIVIDE:       Precedence.PRODUCT,
		MODULO:       Precedence.PRODUCT,
		NEGATE:       Precedence.PREFIX,
		BANG:         Precedence.PREFIX,
		QUESTION:     Precedence.CONDITIONAL,
	}

	precedence, ok := precedences[tok.Type]

	if ok {
		return precedence
	}

	return 0
}
