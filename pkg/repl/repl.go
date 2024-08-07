package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/grantwforsythe/monkeylang/pkg/lexer"
	"github.com/grantwforsythe/monkeylang/pkg/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		// nolint:SA1006
		fmt.Fprintf(out, PROMPT)
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

		_, err := io.WriteString(out, program.String()+"\n")
		if err != nil {
			break
		}
	}
}
