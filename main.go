package main

import (
	"MonkeyInterpreter/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	current, userErr := user.Current()
	if userErr != nil {
		panic(userErr)
	}
	fmt.Printf("Hello %s. Welcome to MonkeyLang.\n", current.Username)
	repl.Start(os.Stdin, os.Stdout)
}