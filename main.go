package main

import (
	"fmt"
	"os/user"

	"github.com/k20ku/monkey/repl"
)

func main() {
	// greeting
	if user, err := user.Current(); err == nil {
		fmt.Printf("Hello %s! ", user.Username)
	} else {
		fmt.Print("Hello! ")
	}
	fmt.Println("This is the Monkey programming language!")
	fmt.Println("Feel free to type in commands")
	fmt.Println("(^D) to exit!")
	fmt.Println()
	// start repl
	repl.Start()
}
