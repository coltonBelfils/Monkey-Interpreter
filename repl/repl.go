package repl

import (
	"MonkeyInterpreter/lexer"
	"MonkeyInterpreter/token"
	"bufio"
	"fmt"
	"io"
)

const (
	PROMPT = ">>"
)

/*
REPL needs to allow for multiline entries. It does this already for multiline comments, but also needs to for things like {} and ()
 */

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	l := lexer.NewChunks()

	fmt.Printf("%s", PROMPT)

	for scanner.Scan() {

		line := scanner.Text()
		l.More(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
		fmt.Printf("%s", PROMPT)
	}
}