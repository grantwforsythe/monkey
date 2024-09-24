// Package token contains all of the tokens for the Monkey language.
package token

type TokenType string

// TODO: Replace assignment operator `let` with `:=`
// TODO: Replace fn defintion with func
// TODO: Add a loop token
// TODO: Add less/greater than or equal to operators
// TODO: Add comment token that will stop evaluation of a line

const (
	ILLEGAL = "ILLEGAL" // Unrecognized token
	EOF     = "EOF"     // End of file

	IDENT  = "IDENT"  // Identifier, e.g. add, foobar, x, y
	INT    = "INT"    // Integer literal, e.g. 1234
	STRING = "STRING" // String literal, "Hello, World!"

	ASSIGN   = "=" // Assignment operator, "="
	PLUS     = "+" // Additional operator, "+"
	MINUS    = "-" // Subtraction operator, "-"
	BANG     = "!" // Boolean inversion operator, "!"
	ASTERISK = "*" // Multiplication operator, "*"
	SLASH    = "/" // Division operator, "/"

	LT = "<" // Less than operator, "<"
	GT = ">" // Greater than operator, ">"

	EQ     = "==" // Equality operator, "=="
	NOT_EQ = "!=" // Inverse equality opertor, "!="

	COMMA     = "," // Comma, ","
	SEMICOLON = ";" // Semicolon, ";"
	COLON     = ":" // Colon, ":"

	LPAREN   = "(" // Left parenthesis, "("
	RPAREN   = ")" // Right parenthesis, ")"
	LBRACE   = "{" // Left brace, "{"
	RBRACE   = "}" // Rigth brace, "}"
	LBRACKET = "[" // Left bracket, "["
	RBRACKET = "]" // Right bracket, "]"

	FUNCTION = "FUNCTION" // Function definition, e.g. "fn(x, y)"
	LET      = "LET"      // Assignment operator, "let"
	TRUE     = "TRUE"     // Boolean literal "true"
	FALSE    = "FALSE"    // Boolean literal "false"
	IF       = "IF"       // Conditonal definition, "if"
	ELSE     = "ELSE"     // Alternative conditional definition, "else"
	RETURN   = "RETURN"   // Return statement, "return"
	MACRO    = "MACRO"    // Macro definition, e.g. "macro(x, y)"
)

type Token struct {
	Type    TokenType // The type of token
	Literal string    // The literal string of the token
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"macro":  MACRO,
}

// Get the token associated with a keyword.
// If the string is not a keyword it is an identifier token.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
