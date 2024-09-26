// Package object contains all of the object definitions.
package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"

	"github.com/grantwforsythe/monkeylang/pkg/ast"
	"github.com/grantwforsythe/monkeylang/pkg/token"
)

type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	QUOTE_OBJ        = "QUOTE"
	MACRO_OBJ        = "MACRO"
)

type Object interface {
	// The type of object
	Type() ObjectType
	// The string representation of the object
	Inspect() string
}

// Represents an Object that is convertible
type Convertible interface {
	// Convert object to resulting AST Node
	ToNode() ast.Node
}

// Represents a hash key for an object
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Represents an object that can be hashed
type Hashable interface {
	// Generate a hash key for an object
	// PERF: Cache return values
	HashKey() HashKey
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}
func (i *Integer) ToNode() ast.Node {
	return &ast.IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(i.Value, 10)},
		Value: i.Value,
	}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}
func (b *Boolean) ToNode() ast.Node {
	var tokenType token.TokenType
	if b.Value {
		tokenType = token.TRUE
	} else {
		tokenType = token.FALSE
	}

	return &ast.BooleanExpression{
		Token: token.Token{Type: tokenType, Literal: strconv.FormatBool(b.Value)},
		Value: b.Value,
	}
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// TODO: Add a stacktrace
// TODO: Add line and column number
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return fmt.Sprintf("Error: %s", e.Message) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range f.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") \n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// A builtin funcion written in the host language (go) and exposed in the interpreter
type BuiltinFunction func(a ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (f *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (f *Builtin) Inspect() string  { return "builtin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range a.Elements {
		elements = append(elements, element.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair // HashPair is the value because we want to keep track of the object used to create the hash
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// An unevaluated AST node
type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType { return QUOTE_OBJ }
func (q *Quote) Inspect() string {
	return fmt.Sprintf("QUOTE(%s)", q.Node.String())
}
func (q *Quote) ToNode() ast.Node {
	return q.Node
}

type Macro struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (m *Macro) Type() ObjectType { return MACRO_OBJ }
func (m *Macro) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range m.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("macro")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n")

	return out.String()
}
