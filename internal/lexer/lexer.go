package lexer

import (
	"unicode"
)

// TokenType defines the type of token.
type TokenType int

const (
	EOF TokenType = iota
	KEYWORD
	STRING_LITERAL
	LPAREN    // (
	RPAREN    // )
	SEMICOLON // ;
	ILLEGAL
)

// Token represents a token with a type and value.
type Token struct {
	Type  TokenType
	Value string
}

// Lexer holds the state of the lexer.
type Lexer struct {
	input  string
	pos    int
	width  int
	tokens []Token
}

// NewLexer initializes a new Lexer.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	for {
		if l.pos >= len(l.input) {
			return Token{Type: EOF}
		}

		ch := l.input[l.pos]
		switch {
		case unicode.IsSpace(rune(ch)):
			l.ignore()
		case isLetter(ch):
			return l.lexKeyword()
		case ch == '"':
			return l.lexStringLiteral()
		case ch == '(':
			l.pos++
			return Token{Type: LPAREN, Value: string(ch)}
		case ch == ')':
			l.pos++
			return Token{Type: RPAREN, Value: string(ch)}
		case ch == ';':
			l.pos++
			return Token{Type: SEMICOLON, Value: string(ch)}
		default:
			l.pos++
			return Token{Type: ILLEGAL, Value: string(ch)}
		}
	}
}

func (l *Lexer) ignore() {
	l.pos++
}

func (l *Lexer) lexKeyword() Token {
	start := l.pos
	for l.pos < len(l.input) && isLetter(l.input[l.pos]) {
		l.pos++
	}
	keyword := l.input[start:l.pos]
	return Token{Type: KEYWORD, Value: keyword}
}

func (l *Lexer) lexStringLiteral() Token {
	start := l.pos
	l.pos++ // skip the opening quote
	for l.pos < len(l.input) && l.input[l.pos] != '"' {
		l.pos++
	}
	l.pos++ // skip the closing quote
	return Token{Type: STRING_LITERAL, Value: l.input[start:l.pos]}
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}
