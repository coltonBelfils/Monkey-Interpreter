This is made by following the *Writing an interpreter in Go* book by Thorsten Ball.

The following additions have been added by me along the way:
- The lexer also encodes in each token the file, line, and position it occurred on. 
- The lexer can receive input from a file, a string, or any io.Reader.
- The lexer works off of Unicode(rune) characters instead of ascii(char). 
- IN PROGRESS: The lexer supports decimals as well as whole numbers