// Package lexer contains a lexer which will tokenize a string.
package lexer

import (
	"github.com/grantwforsythe/monkeylang/pkg/strings"
	"github.com/grantwforsythe/monkeylang/pkg/token"
)

// TODO: Fully support on Unicode and UTF-8 characters (See Chapter 1.3 for more info)

// A structure representing a lexer.
type Lexer struct {
	input        string // The string to be tokenized
	position     int    // current position in input (points to current char)
	readPosition int    // current reading position in input (after current char)
	ch           byte   // current char under examination
}

// Create a new lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Get next character and advance the position in the input string.
// If the current position is greater than the length of the input we've reached the end of the file.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// ASCII "NUL" -> "end of file" or "haven't read anything'"
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

// Peek the next character without advancing the position of the input string.
// If the current position is greater than the length of the input we've reached the end of the file.
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		// EOF
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Skip all the consecutive whitespaces.
func (l *Lexer) skipWhiteSpace() {
	for strings.IsWhiteSpace(l.ch) {
		l.readChar()
	}
}

// Read all consecutive digits.
func (l *Lexer) readDigit() string {
	position := l.position
	for strings.IsDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// Read an identifier and advance the lexer's postions until it encounters a non-letter character.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for strings.IsLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// TODO: Add support for character escaping, e.g. \", \n, etc

// Read the contents of a string.
func (l *Lexer) readString() string {
	// Skip over the first '"'
	l.readChar()

	position := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}

	return l.input[position:l.position]
}

// Iterate to the next token.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
	case ':':
		tok = newToken(token.COLON, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if strings.IsLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// Exit early because we do not want to call readChar twice
			return tok
		} else if strings.IsDigit(l.ch) {
			tok.Literal = l.readDigit()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

// Create a new token
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
