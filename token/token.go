package token

import (
	"fmt"
)

var Precedence = struct {
	LOWEST      int
	ASSIGNMENT  int
	CONDITIONAL int
	COMPARISON  int
	SUM         int
	PRODUCT     int
	EXPONENT    int
	PREFIX      int
	POSTFIX     int
	EXPR        int
}{
	LOWEST:      0,
	ASSIGNMENT:  1,
	CONDITIONAL: 2,
	COMPARISON:  3,
	SUM:         4,
	PRODUCT:     5,
	EXPONENT:    6,
	PREFIX:      7,
	POSTFIX:     8,
	EXPR:        9,
}

type TokenType string

type Token struct {
	Line   int
	Column int
	Value  string
	Type   TokenType
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token{ Line: %d, Column: %d, TokenString: %s, Value: `%s` }",
		t.Line, t.Column, t.Type, t.Value,
	)
}

func (tt TokenType) Str() string {
	var TokenString = map[TokenType]string{
		PLUS:          "+",
		MINUS:         "-",
		NEGATE:        "-",
		BANG:          "!",
		QUESTION:      "?",
		MULTIPLY:      "*",
		DIVIDE:        "/",
		MODULO:        "%",
		EQUAL:         "=",
		LESS:          "<",
		GREATER:       ">",
		LBRACKET:      "[",
		RBRACKET:      "]",
		LBRACE:        "{",
		RBRACE:        "}",
		COLON:         ":",
		LPAREN:        "(",
		RPAREN:        ")",
		COMMA:         ",",
		SEMICOLON:     ";",
		DOT:           ".",
		PLUS_PLUS:     "++",
		MINUS_MINUS:   "--",
		EQUAL_EQUAL:   "==",
		GREATER_EQUAL: ">=",
		LESS_EQUAL:    "<=",
		STRING:        "string",
		INTEGER:       "integer",
		FLOAT:         "float",
		TRUE:          "true",
		FALSE:         "false",
		LET:           "let",
		CONST:         "const",
		IF:            "if",
		ELSE:          "else",
		FOR:           "for",
		SELF:          "self",
		IMPORT:        "import",
		STRUCT:        "struct",
		FUNCTION:      "function",
		IDENTIFIER:    "identifier",
		TYPE:          "types",
		EOL:           "eol",
		EOF:           "eof",
		INVALID:       "invalid",
	}

	return TokenString[tt]
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
	PLUS_PLUS     = "PLUS_PLUS"     // ++
	MINUS_MINUS   = "MINUS_MINUS"   // --
	EQUAL_EQUAL   = "EQUAL_EQUAL"   // ==
	NOT_EQUAL     = "NOT_EQUAL"     // !=
	GREATER_EQUAL = "GREATER_EQUAL" // >=
	LESS_EQUAL    = "LESS_EQUAL"    // <=
	COMMENT       = "COMMENT"       // //

	// Literal
	STRING  = "STRING"
	INTEGER = "INTEGER"
	FLOAT   = "FLOAT"

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
	RETURN     = "RETURN"     // return
	IDENTIFIER = "IDENTIFIER" // any

	// Annotation
	TYPE = "TYPE"

	// Special
	EOL     = "EOL"     // End of line
	EOF     = "EOF"     // End of file
	INVALID = "INVALID" // End of file
)

func Symbol(key string) TokenType {
	keywords := map[string]TokenType{
		"true":   TRUE,
		"false":  FALSE,
		"let":    LET,
		"const":  CONST,
		"if":     IF,
		"else":   ELSE,
		"for":    FOR,
		"self":   SELF,
		"import": IMPORT,
		"struct": STRUCT,
		"fn":     FUNCTION,
		"return": RETURN,
		"int":    TYPE,
		"float":  TYPE,
		"bool":   TYPE,
		"string": TYPE,
	}

	keyword, ok := keywords[key]

	if ok {
		return keyword
	}

	return IDENTIFIER
}

// Precedence is only for infix expression i guess
func GetPrecedence(tok Token) int {
	precedences := map[TokenType]int{
		LBRACE: Precedence.ASSIGNMENT,
		EQUAL:  Precedence.ASSIGNMENT,

		LESS:          Precedence.COMPARISON,
		GREATER:       Precedence.COMPARISON,
		LESS_EQUAL:    Precedence.COMPARISON,
		GREATER_EQUAL: Precedence.COMPARISON,
		EQUAL_EQUAL:   Precedence.COMPARISON,
		NOT_EQUAL:     Precedence.COMPARISON,

		PLUS:  Precedence.SUM,
		MINUS: Precedence.SUM,

		MULTIPLY: Precedence.PRODUCT,
		DIVIDE:   Precedence.PRODUCT,
		MODULO:   Precedence.PRODUCT,

		NEGATE: Precedence.PREFIX,
		BANG:   Precedence.PREFIX,

		QUESTION:    Precedence.CONDITIONAL,
		LPAREN:      Precedence.EXPR,
		LBRACKET:    Precedence.EXPR,
		DOT:         Precedence.EXPR,
		MINUS_MINUS: Precedence.PREFIX,
		PLUS_PLUS:   Precedence.PREFIX,
	}

	precedence, ok := precedences[tok.Type]

	if ok {
		return precedence
	}

	return 0
}
