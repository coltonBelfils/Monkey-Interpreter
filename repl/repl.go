package repl

import (
	"MonkeyInterpreter/evaluator"
	"MonkeyInterpreter/lexer"
	"MonkeyInterpreter/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.NewFromString(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			fmt.Fprintf(out, "%s\n", evaluated.Inspect())
		}

	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(out, "\t%s\n", msg)
	}
}

/*
// I wanted it to look something like this since my lexer natively supports io.Reader,
// but since it requires reading the current token and the next token,
// it gets stuck and can't process the last token until more input has been given.
func Start(in io.Reader, out io.Writer) {
	fmt.Fprintf(out, PROMPT)

	l := lexer.NewFromRepl(in)

	for {
		var tok token.Token
		for tok = l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%v\n", tok)
			if tok.Type == token.QUIT {
				return
			}
		}

		fmt.Fprintf(out, PROMPT)
	}
}
*/
