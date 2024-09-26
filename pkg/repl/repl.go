// Package repl contains the Read Evaluate Print Loop (REPL).
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/grantwforsythe/monkeylang/pkg/evaluator"
	"github.com/grantwforsythe/monkeylang/pkg/lexer"
	"github.com/grantwforsythe/monkeylang/pkg/object"
	"github.com/grantwforsythe/monkeylang/pkg/parser"
)

const PROMPT = ">> "

// Start starts the REPL.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				_, err := io.WriteString(
					out,
					"We ran into some monkey business! Parse errors:\n"+"\t- "+msg.Error()+"\n",
				)
				if err != nil {
					break
				}
			}
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		eval := evaluator.Eval(expanded, env)
		if eval != nil {
			_, err := io.WriteString(out, eval.Inspect()+"\n")
			if err != nil {
				break
			}
		}
	}
}
