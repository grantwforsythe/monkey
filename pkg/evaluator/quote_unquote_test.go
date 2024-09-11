package evaluator

import (
	"testing"

	"github.com/grantwforsythe/monkeylang/pkg/object"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`quote(5)`, `5`},
		{`quote(5 + 8)`, `(5 + 8)`},
		{`quote(foobar)`, `foobar`},
		{`quote(foobar + barfoo)`, `(foobar + barfoo)`},
	}

	for _, test := range tests {
		eval := testEval(test.input)
		quote, ok := eval.(*object.Quote)
		if !ok {
			t.Fatalf("expected *object.Quote. got=%T (%+v)", eval, eval)
		}

		if quote.Node == nil {
			t.Fatalf("quote.Node is nil")
		}

		if quote.Node.String() != test.expected {
			t.Fatalf("quote.Node.String() is not equal to 5. got=%s", quote.Node.String())
		}
	}

}
