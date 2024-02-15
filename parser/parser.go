package parser

import (
	"kat/ast"
	"kat/lexer"
	"kat/token"
	"log"
	"strconv"
)

type PrefixParseFn func() ast.Node
type InfixParseFn func(node ast.Node) ast.Node

type Parser struct {
	Lex       *lexer.Lexer
	Token     token.Token
	PrefixFns map[token.TokenType]PrefixParseFn
	InfixFns  map[token.TokenType]InfixParseFn
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		Lex:       lex,
		PrefixFns: make(map[token.TokenType]PrefixParseFn),
		InfixFns:  make(map[token.TokenType]InfixParseFn),
	}

	// Prefix parsing function
	p.PrefixFns[token.DIGIT] = p.ParseNodeDigit

	// Infix parsing function
	p.InfixFns[token.PLUS] = p.ParseInfix
	p.InfixFns[token.MINUS] = p.ParseInfix
	p.InfixFns[token.MULTIPLY] = p.ParseInfix
	p.InfixFns[token.DIVIDE] = p.ParseInfix
	p.InfixFns[token.MODULO] = p.ParseInfix

	p.NextToken() // prime the token

	return p
}

func (p *Parser) NextToken() token.Token {
	p.Token = p.Lex.NextToken()
	return p.Token
}

func (p *Parser) CurrentToken() token.Token {
	return p.Token
}

func (p *Parser) ParseProgram() ast.NodeProgram {
	program := ast.NodeProgram{}

	for p.CurrentToken().Type != token.EOF {
		if p.CurrentToken().Type == token.EOL {
			p.NextToken()
			continue
		}

		program.Body = append(program.Body, p.ParseStatement())
	}

	return program
}

func (p *Parser) ParseStatement() ast.Node {
	prefixFn, ok := p.PrefixFns[p.CurrentToken().Type]

	if !ok {
		log.Fatalf("PrefixFns parsing function not found for token: %s", p.CurrentToken().Type)
	}

	left := prefixFn()
	p.NextToken()

	for p.CurrentToken().Type != token.EOL && p.CurrentToken().Type != token.EOF {
		infix, ok := p.InfixFns[p.CurrentToken().Type]

		if !ok {
			log.Fatalf("InfixFns parsing function not found for token: %s", p.CurrentToken().Type)
		}

		left = infix(left)
	}

	return left
}

func (p *Parser) ParseNodeDigit() ast.Node {
	val, e := strconv.ParseInt(p.CurrentToken().Value, 10, 64)

	if e != nil {
		log.Fatalf("Parser::Error:%s\n", e)
	}

	return ast.NodeInteger{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseNodeNil() ast.NodeInteger {
	return ast.NodeInteger{
		Token: p.CurrentToken(),
		Value: 999,
	}
}

func (p *Parser) ParseInfix(left ast.Node) ast.Node {
	currentToken := p.CurrentToken()
	p.NextToken()

	right := p.ParseStatement()

	return ast.NodeInfix{
		Token:    currentToken,
		Left:     left,
		Right:    right,
		Operator: currentToken.Value,
	}
}
