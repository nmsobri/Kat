package parser

import (
	"kat/ast"
	"kat/lexer"
	"kat/token"
	"log"
	"strconv"
)

type Parser struct {
	Lex          *lexer.Lexer
	CurrentToken token.Token
}

func New(lex *lexer.Lexer) *Parser {
	return &Parser{Lex: lex}
}

func (p *Parser) Token() token.Token {
	p.CurrentToken = p.Lex.Token()
	return p.CurrentToken
}

func (p *Parser) CCurrentToken() token.Token {
	return p.CurrentToken
}

func (p *Parser) ParseProgram() ast.NodeProgram {
	program := ast.NodeProgram{}

	for p.Token().Type != token.EOF {
		program.Body = append(program.Body, p.ParseStatement())
	}

	return program
}

func (p *Parser) ParseStatement() ast.Node {
	val, e := strconv.ParseInt(p.CCurrentToken().Value, 10, 64)

	if e != nil {
		log.Fatalf("Parser::Error:%s\n", e)
	}

	return ast.NodeInteger{
		Token: p.CCurrentToken(),
		Value: val,
	}
}
