This is made by following the *Writing an interpreter in Go* book by Thorsten Ball.

The following additions have been added by me along the way:
- The lexer also encodes in each token the file, line, and position it occurred on. 
- The lexer can receive input from a file, a string, or any io.Reader.
- The lexer works off of Unicode(rune) characters instead of ascii(char). 
- The lexer supports decimals as well as whole numbers.
- The lexer supports the \, &, |, ^, ~, and % tokens
- The repl can be exited using the quit keyword
- The lexer supports the two character tokens: >=, <=, &&, ||, ++, --, +=, -=, *=, /=, >>, <<
- Numbers are parsed into big.Int and big.Float instead of int.
- The parser supports postfix operators, EG foo++, as well as prefix and infix.

TODO:
- Add comment support to the lexer. Both // and /* */
- Make all errors better conform to the idiomatic Go error format

Might add:
- I want to add a struct data structure to monkey.