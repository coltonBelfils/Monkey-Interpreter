package main

import (
	"MonkeyInterpreter/repl"
	"fmt"
	"os"
	user "os/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hellp %s! Welcom to Monkey!\n", u.Username)

	repl.Start(os.Stdin, os.Stdout)

	fmt.Printf("Goodbye.\n")
}
