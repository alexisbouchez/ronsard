package parser

import (
	"fmt"
	"ronsard/internal/lexer"
)

// NodeType defines the type of a node in the abstract syntax tree (AST).
type NodeType int

const (
	PRINT_STMT NodeType = iota
	STRING_EXPR
)

// Node represents a node in the AST.
type Node struct {
	Type     NodeType
	Value    string
	Children []*Node
}

// Parser holds the state of the parser.
type Parser struct {
	tokens []lexer.Token
	pos    int
}

// NewParser initializes a new Parser.
func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

// Parse parses the tokens and returns the root node of the AST.
func (p *Parser) Parse() (*Node, error) {
	if p.match(lexer.KEYWORD, "afficher") {
		return p.parsePrintStmt()
	}
	return nil, fmt.Errorf("unexpected token: %v", p.currentToken())
}

func (p *Parser) parsePrintStmt() (*Node, error) {
	printNode := &Node{Type: PRINT_STMT}

	// Expect '('
	if !p.match(lexer.LPAREN) {
		return nil, fmt.Errorf("expected '(', got %v", p.currentToken())
	}

	// Expect string literal
	strNode, err := p.parseStringExpr()
	if err != nil {
		return nil, err
	}
	printNode.Children = append(printNode.Children, strNode)

	// Expect ')'
	if !p.match(lexer.RPAREN) {
		return nil, fmt.Errorf("expected ')', got %v", p.currentToken())
	}

	// Expect ';'
	if !p.match(lexer.SEMICOLON) {
		return nil, fmt.Errorf("expected ';', got %v", p.currentToken())
	}

	return printNode, nil
}

func (p *Parser) parseStringExpr() (*Node, error) {
	if p.currentToken().Type == lexer.STRING_LITERAL {
		strNode := &Node{Type: STRING_EXPR, Value: p.currentToken().Value}
		p.pos++
		return strNode, nil
	}
	return nil, fmt.Errorf("expected string literal, got %v", p.currentToken())
}

func (p *Parser) match(expectedType lexer.TokenType, expectedValues ...string) bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	token := p.tokens[p.pos]
	if token.Type != expectedType {
		return false
	}
	if len(expectedValues) > 0 {
		for _, expectedValue := range expectedValues {
			if token.Value == expectedValue {
				p.pos++
				return true
			}
		}
		return false
	}
	p.pos++
	return true
}

func (p *Parser) currentToken() lexer.Token {
	if p.pos >= len(p.tokens) {
		return lexer.Token{Type: lexer.EOF}
	}
	return p.tokens[p.pos]
}
