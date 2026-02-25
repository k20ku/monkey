package main

import (
	"fmt"
	"monkey/repl"
	"os/user"
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

	// start repl
	repl.Start()

	// greeting
	println("Thank you. Goodbye!")
}
